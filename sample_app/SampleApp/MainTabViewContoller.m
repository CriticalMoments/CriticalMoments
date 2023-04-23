//
//  MainTabViewContoller.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-23.
//

#import "MainTabViewContoller.h"

@interface MainTabViewContoller ()

@end

@implementation MainTabViewContoller

- (void)viewDidLoad {
    [super viewDidLoad];
    
    self.view.backgroundColor = [UIColor systemGrayColor];
    
    UIView* bottomLine = [[UIView alloc] init];
    bottomLine.backgroundColor = [UIColor redColor];
    bottomLine.translatesAutoresizingMaskIntoConstraints = NO;
    bottomLine.accessibilityIdentifier = @"bottomBar";
    [self.view addSubview:bottomLine];
    
    UIView* topLine = [[UIView alloc] init];
    topLine.backgroundColor = [UIColor orangeColor];
    topLine.translatesAutoresizingMaskIntoConstraints = NO;
    topLine.accessibilityIdentifier = @"topLine";
    [self.view addSubview:topLine];
    
    NSArray<NSLayoutConstraint*>* constraints = @[
        [bottomLine.heightAnchor constraintEqualToConstant:2.0],
        [bottomLine.leftAnchor constraintEqualToAnchor:self.view.leftAnchor],
        [bottomLine.rightAnchor constraintEqualToAnchor:self.view.rightAnchor],
        [bottomLine.bottomAnchor constraintEqualToAnchor:self.view.layoutMarginsGuide.bottomAnchor],
        [topLine.heightAnchor constraintEqualToConstant:2.0],
        [topLine.leftAnchor constraintEqualToAnchor:self.view.leftAnchor],
        [topLine.rightAnchor constraintEqualToAnchor:self.view.rightAnchor],
        [topLine.topAnchor constraintEqualToAnchor:self.view.layoutMarginsGuide.topAnchor]
    ];
    
    [NSLayoutConstraint activateConstraints:constraints];
}

/*
#pragma mark - Navigation

// In a storyboard-based application, you will often want to do a little preparation before navigation
- (void)prepareForSegue:(UIStoryboardSegue *)segue sender:(id)sender {
    // Get the new view controller using [segue destinationViewController].
    // Pass the selected object to the new view controller.
}
*/

@end
