//
//  MainTabViewContoller.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-23.
//

#import "DemoViewContoller.h"

@interface DemoViewContoller () <UITableViewDataSource, UITableViewDelegate>

@property (nonatomic) CMDemoScreen* screen;

@end

@implementation DemoViewContoller

-(instancetype)initWithDemoScreen:(CMDemoScreen*)screen {
    self = [super init];
    if (self) {
        self.screen = screen;
    }
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];
    
    self.navigationItem.title = @"Critical Moments";
    
    self.view.backgroundColor = [UIColor systemGroupedBackgroundColor];
    
    //self.tableView.dataSource = self;
    
    /*UITableView* tableView = [[UITableView alloc] init];
    [self.view addSubview:tableView];
    
    NSArray<NSLayoutConstraint*>* constraints = @[
        [tableView.leftAnchor constraintEqualToAnchor:self.view.layoutMarginsGuide.leftAnchor],
        [tableView.rightAnchor constraintEqualToAnchor:self.view.layoutMarginsGuide.rightAnchor],
        [tableView.bottomAnchor constraintEqualToAnchor:self.view.layoutMarginsGuide.bottomAnchor],
        [tableView.topAnchor constraintEqualToAnchor:self.view.layoutMarginsGuide.topAnchor],
    ];
    
    [NSLayoutConstraint activateConstraints:constraints];*/
}

-(CMDemoAction*) actionForIndexPath:(NSIndexPath *)indexPath {
    return [self.screen.sections objectAtIndex:indexPath.section].actions[indexPath.row];
}

#pragma mark UITableViewDelegate

-(void)tableView:(UITableView *)tableView didSelectRowAtIndexPath:(NSIndexPath *)indexPath {
    [tableView deselectRowAtIndexPath:indexPath animated:YES];
    CMDemoAction* action = [self actionForIndexPath:indexPath];
    [action performAction];
}

#pragma mark UITableViewDataSource

-(NSInteger)tableView:(UITableView *)tableView numberOfRowsInSection:(NSInteger)section {
    return [self.screen.sections objectAtIndex:section].actions.count;
}

-(NSInteger)numberOfSectionsInTableView:(UITableView *)tableView {
    return self.screen.sections.count;
}

-(UITableViewCell *)tableView:(UITableView *)tableView cellForRowAtIndexPath:(NSIndexPath *)indexPath {
    CMDemoAction* action = [self actionForIndexPath:indexPath];
    
    // TOOD
    //var cell = tableView.dequeueReusableCell(withIdentifier: "myCellType", for: indexPath
    
    UITableViewCell* cell = [[UITableViewCell alloc] initWithStyle:UITableViewCellStyleSubtitle reuseIdentifier:@"demoCell"];
    cell.textLabel.text = action.title;
    cell.detailTextLabel.text = action.subtitle;
    return cell;
}

-(NSString *)tableView:(UITableView *)tableView titleForHeaderInSection:(NSInteger)section {
    return [self.screen.sections objectAtIndex:section].title;
}

@end
