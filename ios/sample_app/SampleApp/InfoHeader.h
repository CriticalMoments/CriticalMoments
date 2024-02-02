//
//  InfoHeader.h
//  SampleApp
//
//  Created by Steve Cosman on 2024-02-02.
//

#import <UIKit/UIKit.h>

#import "CMDemoScreen.h"

NS_ASSUME_NONNULL_BEGIN

@interface InfoHeader : UIView

+ (instancetype)headerWithScreen:(CMDemoScreen *)screen;

@end

NS_ASSUME_NONNULL_END
