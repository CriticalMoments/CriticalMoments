
Pod::Spec.new do |s|
  s.name             = 'CriticalMoments'
  s.version          = '0.1.1-beta'
  s.summary          = 'Deliver the right message, at the right moment.'

  s.description      = <<-DESC
  Deliver the right message, at the right moment.
  
  Work in progress, not ready for use.
                       DESC

  s.homepage         = 'https://criticalmoments.io'
  # s.screenshots     = 'www.example.com/screenshots_1', 'www.example.com/screenshots_2'
  s.license          = { :type => 'CUSTOM', :text => 'Copyright 2023. All rights reserved.' }
  s.author           = { 'Steve Cosman' => 'scosman@users.noreply.github.com' }
  s.source           = { :git => 'https://github.com/criticalmoments/CriticalMoments.git', :tag => s.version.to_s }
  # s.social_media_url = 'https://twitter.com/<TWITTER_USERNAME>'

  s.ios.deployment_target = '11.0'

  s.source_files = 'Sources/**/*'
  s.swift_version = '5.0'
  
  # s.resource_bundles = {
  #   'CriticalMoments' => ['CriticalMoments/Assets/*.png']
  # }

  # s.public_header_files = 'Pod/Classes/**/*.h'
  # s.frameworks = 'UIKit', 'MapKit'
  # s.dependency 'AFNetworking', '~> 2.3'
end

