//
//  MainTabViewContoller.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-23.
//

#import "DemoViewContoller.h"

#import "InfoHeader.h"

#define DEMO_CELL_REUSE_ID @"io.criticalmoments.sample_app.demo_cell"

@interface DemoViewContoller () <UITableViewDataSource, UITableViewDelegate>

@property(nonatomic, strong) CMDemoScreen *screen;
@property(nonatomic, strong) InfoHeader *header;

@end

@implementation DemoViewContoller

- (instancetype)initWithDemoScreen:(CMDemoScreen *)screen {
    self = [super init];
    if (self) {
        self.screen = screen;
        self.header = [InfoHeader headerWithScreen:self.screen];
    }
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];

    self.navigationItem.title = self.screen.title;
    if (@available(iOS 13.0, *)) {
        self.view.backgroundColor = [UIColor systemGroupedBackgroundColor];
    } else {
        self.view.backgroundColor = [UIColor colorWithRed:0.945 green:0.945 blue:0.945 alpha:1.0];
    }
}

- (CMDemoAction *)actionForIndexPath:(NSIndexPath *)indexPath {
    return [[self.screen.sections objectAtIndex:indexPath.section].actions objectAtIndex:indexPath.row];
}

- (void)viewWillLayoutSubviews {
    if (!self.tableView.tableHeaderView) {
        // Delay adding until had broader layout to size
        self.tableView.tableHeaderView = self.header;
    }
    CGSize size = [self.header systemLayoutSizeFittingSize:self.view.frame.size];
    if (self.view.frame.size.width != self.header.frame.size.width || size.height != self.header.frame.size.height) {
        self.header.frame = CGRectMake(0, 0, self.view.frame.size.width, size.height);
        dispatch_async(dispatch_get_main_queue(), ^{
          [self.tableView reloadData];
        });
    }
}

#pragma mark UITableViewDelegate

- (void)tableView:(UITableView *)tableView didSelectRowAtIndexPath:(NSIndexPath *)indexPath {
    [tableView deselectRowAtIndexPath:indexPath animated:YES];
    CMDemoAction *action = [self actionForIndexPath:indexPath];
    [action performAction];
}

#pragma mark UITableViewDataSource

- (NSInteger)tableView:(UITableView *)tableView numberOfRowsInSection:(NSInteger)section {
    return [self.screen.sections objectAtIndex:section].actions.count;
}

- (NSInteger)numberOfSectionsInTableView:(UITableView *)tableView {
    return self.screen.sections.count;
}

- (UITableViewCell *)tableView:(UITableView *)tableView cellForRowAtIndexPath:(NSIndexPath *)indexPath {
    CMDemoAction *action = [self actionForIndexPath:indexPath];

    UITableViewCell *cell = [tableView dequeueReusableCellWithIdentifier:DEMO_CELL_REUSE_ID];
    if (!cell) {
        cell = [[UITableViewCell alloc] initWithStyle:UITableViewCellStyleSubtitle reuseIdentifier:DEMO_CELL_REUSE_ID];
    }
    cell.textLabel.text = action.title;
    cell.detailTextLabel.text = action.subtitle;
    cell.detailTextLabel.numberOfLines = 10;
    cell.hidden = action.skipInUI;
    return cell;
}

- (CGFloat)tableView:(UITableView *)tableView heightForRowAtIndexPath:(NSIndexPath *)indexPath {
    CMDemoAction *action = [self actionForIndexPath:indexPath];
    if (action.skipInUI) {
        return 0;
    }
    return -1; // dynamic
}

- (NSString *)tableView:(UITableView *)tableView titleForHeaderInSection:(NSInteger)section {
    return [self.screen.sections objectAtIndex:section].title;
}

@end
