//
//  MainTabViewContoller.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-23.
//

#import "DemoViewContoller.h"

#define DEMO_CELL_REUSE_ID @"io.criticalmoments.sample_app.demo_cell"

@interface DemoViewContoller () <UITableViewDataSource, UITableViewDelegate>

@property(nonatomic) CMDemoScreen *screen;

@end

@implementation DemoViewContoller

- (instancetype)initWithDemoScreen:(CMDemoScreen *)screen {
    self = [super init];
    if (self) {
        self.screen = screen;
    }
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];

    self.navigationItem.title = self.screen.title;
    self.view.backgroundColor = [UIColor systemGroupedBackgroundColor];
}

- (CMDemoAction *)actionForIndexPath:(NSIndexPath *)indexPath {
    return [[self.screen.sections objectAtIndex:indexPath.section].actions
        objectAtIndex:indexPath.row];
}

#pragma mark UITableViewDelegate

- (void)tableView:(UITableView *)tableView
    didSelectRowAtIndexPath:(NSIndexPath *)indexPath {
    [tableView deselectRowAtIndexPath:indexPath animated:YES];
    CMDemoAction *action = [self actionForIndexPath:indexPath];
    [action performAction];
}

#pragma mark UITableViewDataSource

- (NSInteger)tableView:(UITableView *)tableView
    numberOfRowsInSection:(NSInteger)section {
    return [self.screen.sections objectAtIndex:section].actions.count;
}

- (NSInteger)numberOfSectionsInTableView:(UITableView *)tableView {
    return self.screen.sections.count;
}

- (UITableViewCell *)tableView:(UITableView *)tableView
         cellForRowAtIndexPath:(NSIndexPath *)indexPath {
    CMDemoAction *action = [self actionForIndexPath:indexPath];

    UITableViewCell *cell =
        [tableView dequeueReusableCellWithIdentifier:DEMO_CELL_REUSE_ID];
    if (!cell) {
        cell =
            [[UITableViewCell alloc] initWithStyle:UITableViewCellStyleSubtitle
                                   reuseIdentifier:DEMO_CELL_REUSE_ID];
    }
    cell.textLabel.text = action.title;
    cell.detailTextLabel.text = action.subtitle;
    cell.detailTextLabel.numberOfLines = 4;
    return cell;
}

- (NSString *)tableView:(UITableView *)tableView
    titleForHeaderInSection:(NSInteger)section {
    return [self.screen.sections objectAtIndex:section].title;
}

@end
