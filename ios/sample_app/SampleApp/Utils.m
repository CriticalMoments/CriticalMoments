//
//  Utils.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-26.
//

#import "Utils.h"

@implementation Utils

+ (UIWindow *)keyWindow {
    UIWindow *keyWindow =
        [[[UIApplication sharedApplication] windows] firstObject];
    for (UIWindow *w in [[UIApplication sharedApplication] windows]) {
        if (w.isKeyWindow) {
            keyWindow = w;
            break;
        }
    }
    return keyWindow;
}

@end
