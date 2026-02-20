APP_NAME := sussurro
BUILD_DIR := bin
CMD_DIR := cmd/sussurro

# Whisper.cpp configuration
WHISPER_DIR := third_party/whisper.cpp
WHISPER_INCLUDE := $(abspath $(WHISPER_DIR)/include)
WHISPER_GGML_INCLUDE := $(abspath $(WHISPER_DIR)/ggml/include)
C_INCLUDE_PATH := $(WHISPER_INCLUDE):$(WHISPER_GGML_INCLUDE)
LIBRARY_PATH := $(abspath $(WHISPER_DIR))

# go-llama.cpp configuration
LLAMA_DIR := third_party/go-llama.cpp

# Detect number of CPU cores for parallel builds
NPROCS := $(shell nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 1)

# Detect OS for platform-specific builds
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	BUILD_TYPE := metal
	GGML_METAL_PATH := -L$(WHISPER_DIR)/build/ggml/src/ggml-metal
else
	BUILD_TYPE :=
	GGML_METAL_PATH :=
endif

# ---- UI / overlay dependencies (Linux only) ----
HAS_LAYER_SHELL    := $(shell pkg-config --exists gtk-layer-shell          2>/dev/null && echo yes || echo no)
HAS_AYATANA        := $(shell pkg-config --exists ayatana-appindicator3-0.1 2>/dev/null && echo yes || echo no)
HAS_APPINDICATOR   := $(shell pkg-config --exists appindicator3-0.1         2>/dev/null && echo yes || echo no)

LAYER_CFLAGS  := $(shell pkg-config --cflags gtk+-3.0 2>/dev/null)
LAYER_LDFLAGS := $(shell pkg-config --libs   gtk+-3.0 2>/dev/null)

ifeq ($(HAS_LAYER_SHELL),yes)
LAYER_CFLAGS  += $(shell pkg-config --cflags gtk-layer-shell 2>/dev/null) -DHAVE_GTK_LAYER_SHELL
LAYER_LDFLAGS += $(shell pkg-config --libs   gtk-layer-shell 2>/dev/null)
endif

WV_CFLAGS  := $(shell pkg-config --cflags webkit2gtk-4.1 2>/dev/null || pkg-config --cflags webkit2gtk-4.0 2>/dev/null)
WV_LDFLAGS := $(shell pkg-config --libs   webkit2gtk-4.1 2>/dev/null || pkg-config --libs   webkit2gtk-4.0 2>/dev/null)

# If only webkit2gtk-4.1 is available, create a compat .pc file so that
# webview_go (which hardcodes pkg-config: webkit2gtk-4.0) can find it.
HAS_WV40 := $(shell pkg-config --exists webkit2gtk-4.0 2>/dev/null && echo yes || echo no)
HAS_WV41 := $(shell pkg-config --exists webkit2gtk-4.1 2>/dev/null && echo yes || echo no)

ifeq ($(HAS_WV40),no)
ifeq ($(HAS_WV41),yes)
COMPAT_PC_DIR := $(abspath .build-compat/pkgconfig)
PKG_CONFIG_PATH_UI := $(COMPAT_PC_DIR)$(if $(PKG_CONFIG_PATH),:$(PKG_CONFIG_PATH),)
else
$(warning Neither webkit2gtk-4.0 nor webkit2gtk-4.1 found; UI build will fail)
COMPAT_PC_DIR :=
PKG_CONFIG_PATH_UI :=
endif
else
COMPAT_PC_DIR :=
PKG_CONFIG_PATH_UI := $(PKG_CONFIG_PATH)
endif

# Build tags: use legacy_appindicator when ayatana is not available but appindicator3 is
UI_TAGS :=
ifeq ($(UNAME_S),Linux)
ifeq ($(HAS_AYATANA),no)
ifeq ($(HAS_APPINDICATOR),yes)
UI_TAGS := -tags legacy_appindicator
endif
endif
endif

# Base CGO link flags (whisper + llama)
BASE_LDFLAGS := -L$(WHISPER_DIR)/build/src -L$(WHISPER_DIR)/build/ggml/src \
	-L$(WHISPER_DIR)/build/ggml/src/ggml-cpu $(GGML_METAL_PATH) \
	-L$(WHISPER_DIR)/build/ggml/src/ggml-blas -lwhisper

# Export environment variables for CGO
export C_INCLUDE_PATH
export LIBRARY_PATH

.PHONY: all build compat-pc run clean deps

all: build

deps:
	@mkdir -p third_party
	@if [ ! -d "$(WHISPER_DIR)" ]; then \
		echo "Cloning whisper.cpp..."; \
		git clone https://github.com/ggerganov/whisper.cpp.git $(WHISPER_DIR); \
		echo "Patching whisper.cpp symbols..."; \
		chmod +x scripts/patch-whisper.sh; \
		./scripts/patch-whisper.sh; \
	fi
	@echo "Building whisper.cpp library..."
	@cmake -S $(WHISPER_DIR) -B $(WHISPER_DIR)/build -DGGML_NATIVE=OFF -DBUILD_SHARED_LIBS=OFF -DWHISPER_BUILD_TESTS=OFF -DWHISPER_BUILD_EXAMPLES=OFF
	@cmake --build $(WHISPER_DIR)/build --config Release --target whisper -j $(NPROCS)
	@if [ ! -d "$(LLAMA_DIR)" ]; then \
		echo "Cloning go-llama.cpp..."; \
		git clone --recursive https://github.com/AshkanYarmoradi/go-llama.cpp $(LLAMA_DIR); \
	fi
	@echo "Building go-llama.cpp library..."
	@$(MAKE) -C $(LLAMA_DIR) clean
	@$(MAKE) -j $(NPROCS) -C $(LLAMA_DIR) libbinding.a BUILD_TYPE=$(BUILD_TYPE)

# Create webkit2gtk-4.0 compatibility .pc when only 4.1 is installed
compat-pc:
ifneq ($(COMPAT_PC_DIR),)
	@mkdir -p $(COMPAT_PC_DIR)
	@printf 'Name: webkit2gtk-4.0\nDescription: WebKit2 GTK+ (4.1 compat)\nVersion: 2.99.0\nRequires: webkit2gtk-4.1\nLibs: %s\nCflags: %s\n' \
		"$(shell pkg-config --libs webkit2gtk-4.1)" \
		"$(shell pkg-config --cflags webkit2gtk-4.1)" \
		> $(COMPAT_PC_DIR)/webkit2gtk-4.0.pc
	@echo "  Created compat .pc: $(COMPAT_PC_DIR)/webkit2gtk-4.0.pc"
endif

# Build with full UI (overlay + tray + settings window)
build: deps compat-pc
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
ifeq ($(UNAME_S),Darwin)
	CGO_LDFLAGS="$(BASE_LDFLAGS) -framework Cocoa -framework QuartzCore -framework CoreVideo -framework Foundation" \
	go build -o $(BUILD_DIR)/$(APP_NAME) ./$(CMD_DIR)
else
	@echo "  Layer shell  : $(HAS_LAYER_SHELL)"
	@echo "  Ayatana tray : $(HAS_AYATANA)"
	@echo "  AppIndicator : $(HAS_APPINDICATOR)"
	@echo "  Build tags   : $(UI_TAGS)"
	PKG_CONFIG_PATH="$(PKG_CONFIG_PATH_UI)" \
	CGO_CFLAGS="$(LAYER_CFLAGS) $(WV_CFLAGS)" \
	CGO_LDFLAGS="$(BASE_LDFLAGS) $(LAYER_LDFLAGS) $(WV_LDFLAGS)" \
	go build $(UI_TAGS) -o $(BUILD_DIR)/$(APP_NAME) ./$(CMD_DIR)
endif

run: build
	@echo "Running $(APP_NAME)..."
	@./$(BUILD_DIR)/$(APP_NAME)

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -rf third_party
	@rm -rf .build-compat
