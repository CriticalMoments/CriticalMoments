#
# Be sure to run `pod lib lint CriticalMoments.podspec' to ensure this is a
# valid spec before submitting.
#
# Any lines starting with a # are optional, but their use is encouraged
# To learn more about a Podspec see https://guides.cocoapods.org/syntax/podspec.html
#

Pod::Spec.new do |s|
  s.name             = 'CriticalMoments'
  s.version          = '0.1.0-beta'
  s.summary          = 'Deliver the right message, at the right moment.'

# This description is used to generate tags and improve search results.
#   * Think: What does it do? Why did you write it? What is the focus?
#   * Try to keep it short, snappy and to the point.
#   * Write the description between the DESC delimiters below.
#   * Finally, don't worry about the indent, CocoaPods strips it!

  s.description      = <<-DESC
  Deliver the right message, at the right moment.
  
  Work in progress, not ready for use.
                       DESC

  s.homepage         = 'https://criticalmoments.io'
  # s.screenshots     = 'www.example.com/screenshots_1', 'www.example.com/screenshots_2'
  s.license          = { :type => 'CUSTOM', :text => 'Copyright 2023. All rights reserved.' }
  s.author           = { 'scosman' => 'scosman@users.noreply.github.com' }
  s.source           = { :git => 'https://github.com/criticalmoments/CriticalMoments.git', :tag => s.version.to_s }
  # s.social_media_url = 'https://twitter.com/<TWITTER_USERNAME>'

  s.ios.deployment_target = '10.0'

  s.source_files = 'ios/CriticalMoments/Classes/**/*'
  
  # s.resource_bundles = {
  #   'CriticalMoments' => ['CriticalMoments/Assets/*.png']
  # }

  # s.public_header_files = 'Pod/Classes/**/*.h'
  # s.frameworks = 'UIKit', 'MapKit'
  # s.dependency 'AFNetworking', '~> 2.3'
end
