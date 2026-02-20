// overlay_darwin.m — NSPanel overlay with CoreGraphics drawing
#import <Cocoa/Cocoa.h>
#import <QuartzCore/QuartzCore.h>
#include <math.h>

/* ---- State constants (must match Go) ---- */
#define OVERLAY_STATE_IDLE          0
#define OVERLAY_STATE_RECORDING     1
#define OVERLAY_STATE_TRANSCRIBING  2

#define ITEM_COUNT     7
#define BAR_MIN_HEIGHT 4.0
#define BAR_MAX_HEIGHT 40.0
#define RMS_SCALE      0.08

typedef void (*HotkeyDownCB)(void);
typedef void (*HotkeyUpCB)(void);

/* ------------------------------------------------------------------ */
/* SussurroView — NSView subclass that draws the pill                  */
/* ------------------------------------------------------------------ */

@interface SussurroView : NSView {
@public
    int     state;
    double  animTime;
    double  shimmerPhase;

    float   rmsRing[7];
    int     rmsHead;
    double  barHeights[7];
    double  barTargets[7];

    CVDisplayLinkRef displayLink;
}
- (void)tick:(double)dt;
@end

static CVReturn displayLinkCallback(CVDisplayLinkRef link,
                                    const CVTimeStamp *now,
                                    const CVTimeStamp *output,
                                    CVOptionFlags flagsIn,
                                    CVOptionFlags *flagsOut,
                                    void *ctx)
{
    (void)link; (void)now; (void)output; (void)flagsIn; (void)flagsOut;
    SussurroView *v = (__bridge SussurroView *)ctx;
    dispatch_async(dispatch_get_main_queue(), ^{
        [v tick:1.0/60.0];
    });
    return kCVReturnSuccess;
}

@implementation SussurroView

- (instancetype)initWithFrame:(NSRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        state        = OVERLAY_STATE_IDLE;
        animTime     = 0.0;
        shimmerPhase = 0.0;
        rmsHead      = 0;
        for (int i = 0; i < ITEM_COUNT; i++) {
            barHeights[i] = BAR_MIN_HEIGHT;
            barTargets[i] = BAR_MIN_HEIGHT;
        }

        CVDisplayLinkCreateWithActiveCGDisplays(&displayLink);
        CVDisplayLinkSetOutputCallback(displayLink, displayLinkCallback,
                                       (__bridge void *)self);
        CVDisplayLinkStart(displayLink);
    }
    return self;
}

- (void)dealloc
{
    CVDisplayLinkStop(displayLink);
    CVDisplayLinkRelease(displayLink);
}

- (void)tick:(double)dt
{
    animTime     += dt;
    shimmerPhase += dt;
    for (int i = 0; i < ITEM_COUNT; i++) {
        barHeights[i] = barHeights[i] * 0.7 + barTargets[i] * 0.3;
    }
    [self setNeedsDisplay:YES];
}

- (BOOL)isOpaque { return NO; }
- (BOOL)wantsLayer { return YES; }

- (void)drawRect:(NSRect)dirtyRect
{
    (void)dirtyRect;

    CGContextRef ctx = [[NSGraphicsContext currentContext] CGContext];
    NSRect bounds = self.bounds;
    double w = bounds.size.width;
    double h = bounds.size.height;
    double r = h / 2.0;  /* radius = half height for true pill */

    /* Clear */
    CGContextClearRect(ctx, bounds);

    /* Pill background */
    CGMutablePathRef path = CGPathCreateMutable();
    CGPathAddRoundedRect(path, NULL,
                         CGRectMake(0, 0, w, h), r, r);
    CGContextAddPath(ctx, path);
    CGPathRelease(path);
    CGContextSetRGBFillColor(ctx, 0.102, 0.102, 0.102, 0.90);
    CGContextFillPath(ctx);

    switch (state) {
    case OVERLAY_STATE_IDLE:   [self drawDots:ctx w:w h:h]; break;
    case OVERLAY_STATE_RECORDING:  [self drawBars:ctx w:w h:h]; break;
    case OVERLAY_STATE_TRANSCRIBING: [self drawShimmer:ctx w:w h:h]; break;
    }
}

- (void)drawDots:(CGContextRef)ctx w:(double)w h:(double)h
{
    double spacing  = 10.0, dotR = 3.0;
    double totalW   = (ITEM_COUNT - 1) * spacing;
    double startX   = (w - totalW) / 2.0;
    double cy       = h / 2.0;

    for (int i = 0; i < ITEM_COUNT; i++) {
        double phi = 2.0 * M_PI * animTime / 4.0 + i * 2.0 * M_PI / (double)ITEM_COUNT;
        double s   = sin(phi);
        double a   = 0.35 + 0.65 * s * s;
        double cx  = startX + i * spacing;
        CGContextSetRGBFillColor(ctx, 1, 1, 1, a);
        CGContextFillEllipseInRect(ctx,
            CGRectMake(cx - dotR, cy - dotR, dotR*2, dotR*2));
    }
}

- (void)drawBars:(CGContextRef)ctx w:(double)w h:(double)h
{
    double spacing = 8.0, bw = 5.0, br = 2.5;
    double totalW  = (ITEM_COUNT - 1) * spacing;
    double startX  = (w - totalW) / 2.0;
    double cy      = h / 2.0;

    CGContextSetRGBFillColor(ctx, 1, 1, 1, 1);
    for (int i = 0; i < ITEM_COUNT; i++) {
        double bh = barHeights[i];
        double cx = startX + i * spacing;
        double x  = cx - bw / 2.0;
        double y  = cy - bh / 2.0;
        CGRect  rect = CGRectMake(x, y, bw, bh);
        CGFloat ri = br;
        CGMutablePathRef p = CGPathCreateMutable();
        CGPathAddRoundedRect(p, NULL, rect, ri, ri);
        CGContextAddPath(ctx, p);
        CGPathRelease(p);
        CGContextFillPath(ctx);
    }
}

- (void)drawShimmer:(CGContextRef)ctx w:(double)w h:(double)h
{
    /* Draw "transcribing" text with a moving shimmer */
    CGContextSetRGBFillColor(ctx, 1, 1, 1, 0.7);

    NSAttributedString *str = [[NSAttributedString alloc]
        initWithString:@"transcribing"
            attributes:@{
                NSFontAttributeName: [NSFont systemFontOfSize:14],
                NSForegroundColorAttributeName: [NSColor colorWithWhite:1 alpha:0.7]
            }];

    NSSize sz = str.size;
    NSPoint pt = NSMakePoint((w - sz.width) / 2.0,
                             (h - sz.height) / 2.0);
    [str drawAtPoint:pt];

    /* Shimmer gradient overlay */
    double phase    = fmod(shimmerPhase, 1.5) / 1.5;
    double sx       = pt.x - 40.0 + (sz.width + 80.0) * phase;
    NSGradient *grad = [[NSGradient alloc]
        initWithColors:@[
            [NSColor colorWithWhite:1 alpha:0.0],
            [NSColor colorWithWhite:1 alpha:0.5],
            [NSColor colorWithWhite:1 alpha:0.0]
        ]
        atLocations:(CGFloat[]){0.0, 0.5, 1.0}
        colorSpace:[NSColorSpace genericRGBColorSpace]];

    NSRect shimRect = NSMakeRect(sx - 20, pt.y, 40, sz.height);
    [grad drawInRect:shimRect angle:0];
}

@end

/* ------------------------------------------------------------------ */
/* SussurroPanel — NSPanel subclass                                    */
/* ------------------------------------------------------------------ */

@interface SussurroPanel : NSPanel
@end
@implementation SussurroPanel
- (BOOL)canBecomeKeyWindow { return NO; }
- (BOOL)canBecomeMainWindow { return NO; }
@end

/* ------------------------------------------------------------------ */
/* C-linkage API                                                       */
/* ------------------------------------------------------------------ */

static SussurroPanel *g_panel = nil;
static SussurroView  *g_view  = nil;

void* overlay_create_macos(void)
{
    NSScreen *screen = [NSScreen mainScreen];
    NSRect    sf     = screen.frame;
    NSRect    frame  = NSMakeRect(
        (sf.size.width  - 220) / 2.0,
        24,
        220, 52);

    g_panel = [[SussurroPanel alloc]
        initWithContentRect:frame
                  styleMask:NSWindowStyleMaskBorderless
                    backing:NSBackingStoreBuffered
                      defer:NO];

    g_panel.level             = NSFloatingWindowLevel + 1;
    g_panel.opaque            = NO;
    g_panel.hasShadow         = NO;
    g_panel.backgroundColor   = [NSColor clearColor];
    g_panel.collectionBehavior =
        NSWindowCollectionBehaviorCanJoinAllSpaces |
        NSWindowCollectionBehaviorStationary |
        NSWindowCollectionBehaviorIgnoresCycle;

    g_view = [[SussurroView alloc] initWithFrame:
              NSMakeRect(0, 0, frame.size.width, frame.size.height)];
    [g_panel setContentView:g_view];
    [g_panel makeKeyAndOrderFront:nil];

    return (__bridge void *)g_panel;
}

void overlay_set_state_macos(int state)
{
    dispatch_async(dispatch_get_main_queue(), ^{
        if (g_view) g_view->state = state;
    });
}

void overlay_push_rms_macos(float rms)
{
    dispatch_async(dispatch_get_main_queue(), ^{
        if (!g_view) return;
        g_view->rmsRing[g_view->rmsHead] = rms;
        g_view->rmsHead = (g_view->rmsHead + 1) % ITEM_COUNT;
        for (int i = 0; i < ITEM_COUNT; i++) {
            int idx = (g_view->rmsHead + i) % ITEM_COUNT;
            float v = g_view->rmsRing[idx];
            double norm = v / RMS_SCALE;
            if (norm > 1.0) norm = 1.0;
            g_view->barTargets[i] = BAR_MIN_HEIGHT +
                                    norm * (BAR_MAX_HEIGHT - BAR_MIN_HEIGHT);
        }
    });
}

void overlay_show_macos(void)
{
    dispatch_async(dispatch_get_main_queue(), ^{
        [g_panel makeKeyAndOrderFront:nil];
    });
}

void overlay_hide_macos(void)
{
    dispatch_async(dispatch_get_main_queue(), ^{
        [g_panel orderOut:nil];
    });
}
