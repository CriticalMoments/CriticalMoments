//
//  CMNotificationObserver.m
//
//
//  Created by Steve Cosman on 2024-02-05.
//

#import "CMNotificationObserver.h"

#import "../CriticalMoments_private.h"

@import Appcore;

@interface CMNotificationObserver ()

@property(nonatomic, weak) CriticalMoments *cm;
@property(nonatomic) bool starting, started;

@end

@implementation CMNotificationObserver

- (instancetype)initWithCm:(CriticalMoments *)cm {
    self = [super init];
    if (self) {
        self.cm = cm;
    }
    return self;
}

- (void)start {
    @synchronized(self) {
        if (self.started || self.starting) {
            return;
        }

        self.starting = true;
    }

    // register observers
    [[NSNotificationCenter defaultCenter] addObserver:self
                                             selector:@selector(processNotification:)
                                                 name:UIApplicationDidEnterBackgroundNotification
                                               object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self
                                             selector:@selector(processNotification:)
                                                 name:UIApplicationWillEnterForegroundNotification
                                               object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self
                                             selector:@selector(processNotification:)
                                                 name:UIApplicationWillTerminateNotification
                                               object:nil];

    // get current state (needs to be in main queue)
    dispatch_async(dispatch_get_main_queue(), ^{
      UIApplicationState state = [[UIApplication sharedApplication] applicationState];
      if (state == UIApplicationStateBackground) {
          [self.cm sendEvent:DatamodelAppEnteredBackgroundBuiltInEvent builtIn:true handler:nil];
      } else {
          // Forground == inactive or active states
          [self.cm sendEvent:DatamodelAppEnteredForegroundBuiltInEvent builtIn:true handler:nil];
      }
      self.started = true;
    });
}

- (void)processNotification:(NSNotification *)notification {
    if (!self.started) {
        // ignore notifications until started. The first update should come from the startup code checking state
        return;
    }

    if (UIApplicationDidEnterBackgroundNotification == notification.name) {
        [self.cm sendEvent:DatamodelAppEnteredBackgroundBuiltInEvent builtIn:true handler:nil];
    } else if (UIApplicationWillEnterForegroundNotification == notification.name) {
        [self.cm sendEvent:DatamodelAppEnteredForegroundBuiltInEvent builtIn:true handler:nil];
    } else if (UIApplicationWillTerminateNotification == notification.name) {
        [self.cm sendEvent:DatamodelAppTerminatedBuiltInEvent builtIn:true handler:nil];
    }
}

@end
