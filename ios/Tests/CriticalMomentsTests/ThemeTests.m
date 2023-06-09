//
//  ThemeTests.m
//
//
//  Created by Steve Cosman on 2023-05-03.
//

#import <XCTest/XCTest.h>

#import "themes/CMTheme.h"
#import "themes/CMTheme_private.h"

@interface ThemeTests : XCTestCase

@end

@implementation ThemeTests

- (void)setUp {
    // Put setup code here. This method is called before the invocation of each
    // test method in the class.
}

- (void)tearDown {
    // Put teardown code here. This method is called after the invocation of
    // each test method in the class.
}

// TODO remove test theme, use E2E config named/default theme instead
- (void)testAppcoreThemeEndToEnd {
    CMTheme *testTheme = [CMTheme testTheme];

    // banners
    XCTAssert([UIColor.greenColor isEqual:testTheme.bannerForegroundColor],
              @"banner foreground should be green on test theme");
    XCTAssert([UIColor.redColor isEqual:testTheme.bannerBackgroundColor],
              @"banner background should be red on test theme");

    // colors
    XCTAssert([UIColor.redColor isEqual:[testTheme primaryColorForView:[[UIView alloc] init]]],
              @"primary color should be red on test theme");
    UIColor *white = [UIColor colorWithRed:1 green:1 blue:1 alpha:1];
    XCTAssert([white isEqual:testTheme.backgroundColor], @"background should be white on test theme");
    XCTAssert([UIColor.redColor isEqual:testTheme.primaryTextColor], @"primary text should be red on test theme");
    XCTAssert([UIColor.greenColor isEqual:testTheme.secondaryTextColor],
              @"secondary text should be green on test theme");

    // fonts
    XCTAssert(fabs(testTheme.fontScale - 1.1) < FLT_EPSILON, @"font scale integration issue");
    XCTAssert(!testTheme.scaleFontForDynamicType, @"dynamic type integration issue");
    XCTAssert([@"Palatino-Roman" isEqualToString:testTheme.fontName], @"font name mismatch");
    XCTAssert([@"Palatino-Bold" isEqualToString:testTheme.boldFontName], @"font name mismatch");

    // Dark mode reverses colors
    XCTAssert([UIColor.redColor isEqual:testTheme.darkModeTheme.bannerForegroundColor],
              @"banner foreground should be green on test theme");
    XCTAssert([UIColor.greenColor isEqual:testTheme.darkModeTheme.bannerBackgroundColor],
              @"banner background should be red on test theme");
}

@end
