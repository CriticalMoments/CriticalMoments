
<p align="center">
  <a href="https://criticalmoments.io">
    <img width="320" alt="Critical Moments Logo with 'The Mobile Growth SDK' tagline" src="https://github.com/CriticalMoments/CriticalMoments/assets/848343/9f985505-264b-4b61-af7c-e79f15d01d54">
  </a>
</p>

<p align="center">
  <a href="https://github.com/CriticalMoments/CriticalMoments/actions/workflows/test_release.yml" target="_blank"><img src="https://github.com/CriticalMoments/CriticalMoments/actions/workflows/test_release.yml/badge.svg" alt="Release Tests"></a>
  <a href="https://github.com/CriticalMoments/CriticalMoments/blob/main/test_count.sh"><img src="https://img.shields.io/badge/Test_Case_Count-2550-brightgreen?logo=github&labelColor=32383f&logoColor=969da4" alt="Test Case Count" /></a>
  <a href="https://github.com/CriticalMoments/CriticalMoments/releases/latest"><img src="https://img.shields.io/github/v/release/CriticalMoments/CriticalMoments?color=brightgreen&labelColor=32383f&label=SPM%20Release" alt="Test Case Count" /></a>
</p>

<p align="center">
  <a href="https://docs.criticalmoments.io/quick-start"><strong>Quick Start</strong></a> •
  <a href="https://criticalmoments.io"><strong>Homepage</strong></a> •
  <a href="https://docs.criticalmoments.io"><strong>Documentation</strong></a> • 
  <a href="https://github.com/CriticalMoments/CriticalMoments/issues"><strong>Issues</strong></a>
</p>


# Critical Moments

The **Mobile Growth** SDK. 

We’ve designed the [Critical Moments SDK](https://criticalmoments.io) to be the ultimate growth tool for iOS apps. Our goal is to automate tedious and repetitive growth tasks while providing powerful new growth tools.


## Table of Contents
1. [How It Works](#how-it-works-)
2. [What Makes Critical Moments Special](#what-makes-critical-moments-special)
3. [Quick Start](#quick-start-)
4. [Demo App](#demo-app-)
5. [Homepage, Docs, License & Copyright](#homepage-)

## How It Works 👩‍💻

### Step 1: Install our SDK

Install our SDK and integrate into your app following our [Quick-Start Guide](https://docs.criticalmoments.io/quick-start). This only takes about 15 minutes.

### Step 2: Create your growth plan in JSON

Yes, really—a growth plan defined in JSON. Defining your growth plan in a config file allows you to update it over the air without app updates, deploy new growth tools without new code, change targeting logic anytime, and easily use our [growing library](https://criticalmoments.io/blog) of pre-defined and proven growth tactics. 

- **Add notifications**: we have templates to [increase activation](https://docs.criticalmoments.io/guides/reduce-app-churn-with-notifications#increase-activation-rate), [reduce churn](https://docs.criticalmoments.io/guides/reduce-app-churn-with-notifications#reduce-long-term-churn), and [custom notifications](https://docs.criticalmoments.io/guides/reduce-app-churn-with-notifications#step-5-add-custom-notification) 
- **Improve your app rating**: use our template to [ask users to rate you at the perfect moment](https://docs.criticalmoments.io/guides/improve-your-app-store-rating), increasing your rating and rating volume.
- **Optimize revenue**: ask users to [upgrade at the perfect moment](https://criticalmoments.io/features/grow_revenue), with [over 100 targeting properties built-in](https://docs.criticalmoments.io/conditional-targeting/built-in-properties)
- **In-app messaging**: add [fully native messaging UI](https://docs.criticalmoments.io/actions-in-app-messaging/actions-overview) without writing any code. Make announcements with [banners](https://docs.criticalmoments.io/actions-in-app-messaging/banners), [modals](https://docs.criticalmoments.io/actions-in-app-messaging/modals), [alerts](https://docs.criticalmoments.io/actions-in-app-messaging/alerts), [browser](https://docs.criticalmoments.io/actions-in-app-messaging/open-link), and more. All [themed](https://docs.criticalmoments.io/themes/theme-overview) to match your brand.
- **Smart feature flags**: define [feature flags](https://docs.criticalmoments.io/guides/feature-flags-guide) that can be updated based on over [100 realtime device conditions \(low battery, has network, permissions, etc\)](https://docs.criticalmoments.io/conditional-targeting/built-in-properties). 

### Step 3: Update Anytime, Without App Updates

Once your initial growth plan is deployed, you can update anytime without app updates. This helps in a number of ways:

- [Disaster Recovery](https://criticalmoments.io/features/disaster_recovery): Quickly recover from unexpected bugs, outages, deprecations, and other critical events without negative reviews
- Deploy new growth tactics without additional engineering work, and outside app-release cycles
- Make targeted announcements like pricing changes and legal updates with [in-app messaging](https://criticalmoments.io/features/in_app_messaging)
- Encourage any users on outdated or buggy releases to update their app

## What Makes Critical Moments Special?

### Powerful Targeting 🎯

With Critical Moments, you can target users with the right actions at the perfect moment.

Our simple [string-based conditional statements](https://docs.criticalmoments.io/conditional-targeting/intro-to-conditions) can check over [100 built-in properties](https://docs.criticalmoments.io/conditional-targeting/built-in-properties), custom properties, in-app events, and user engagement history. 

Some examples: 
- `eventCount('app_launch') > 5 && latestEventTime('asked_to_subscribe') < now() - duration('72h')`
- `device_model_class == 'iPad' && versionLessThan(app_version, '2.4.1')`
- `camera_permission != 'authorized' && photo_library_permission != 'authorized'`
- `has_watch || location_city == 'Toronto' || has_car_audio || on_call || has_bt_headset || network_connection_type == 'cellular'`
- `weather_condition IN ['Rain', 'Thunderstorms'] || weather_cloud_cover > 0.80`

### User Privacy 🔑🔒

Critical Moments is designed from the ground up for user privacy. All logic is run locally on their own device. We don’t collect any information about your users. A default installation makes zero calls to our servers from the user’s device*.

*Some optional services like GeoIP location and weather require a service. These services are clearly outlined in our docs, are completely optional, don’t collect user identifiers, and don’t store logs long-term.

### Powerful Config-Driven Growth Plans 📈

Our growth config file format supports building complex growth logic, entirely in config. Connect events in your app to messaging, notifications, paywalls, review prompts, and much more. Check for the perfect moment with conditions. Non-technical team members can contribute, without writing code. Update anytime your config anytime, without App Store updates. 

### Local User-Engagement Database 📙

Our SDK automatically starts building an on-device database of user actions. The most [commonly needed actions are tracked automatically](https://docs.criticalmoments.io/events/built-in-events). [Add your own custom events](https://docs.criticalmoments.io/events/event-overview) and [properties](https://docs.criticalmoments.io/conditional-targeting/custom-properties) in a single line of code. Use this database when targeting user-messaging, reviews, notifications, paywalls and more!

## Quick Start 🚀

Read our [Quick-Start](https://docs.criticalmoments.io/quick-start) guide to get up and running in minutes. 

## Demo App 

Want to see Critical Moments in action? Download our [demo app from TestFlight](https://testflight.apple.com/join/uSwscwu0) or view the source code [on Github](https://github.com/CriticalMoments/CriticalMoments/tree/main/ios/sample_app).

## Homepage 🏠

Check out our [homepage](https://criticalmoments.io) for information about Critical Moments, pricing, and account login.

## Documentation 👩‍💻

Check out our [documentation](https://docs.criticalmoments.io) for details on how to use Critical Moments.

## License

Critical Moments requires you to purchase a license to use it in production apps. See our [pricing page](https://criticalmoments.io/pricing) for details. The code of the SDK is fully source-available, and in this repo.

## Copyright

Copyright (c) 2023 Chesterfield Laboratories Inc.

"Critical Moments" and our logos are trademarks. All rights reserved.
