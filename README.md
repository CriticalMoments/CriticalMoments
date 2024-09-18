
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

Our **Mobile Growth SDK** is designed to: 

> 1) Automate tedious and repetitive growth tasks.
> 2) Provide powerful new growth tools.

## Overview 🔭

- **Growth plans defined in JSON**: Push updates anytime without app store reviews. Make updates without writing new code. Proven templates to get you started.
- **Rich Growth Features**: [smart notifications](https://docs.criticalmoments.io/guides/reduce-app-churn-with-notifications), [app-reviews](https://docs.criticalmoments.io/guides/improve-your-app-store-rating), [paywall timing](https://criticalmoments.io/features/grow_revenue), [native modal UI](https://docs.criticalmoments.io/actions-in-app-messaging/modals), [app-wide banners](https://docs.criticalmoments.io/actions-in-app-messaging/banners), [smart-feature flags](https://docs.criticalmoments.io/guides/feature-flags-guide), and [more](https://docs.criticalmoments.io/concepts-overview). 
- **Powerful Targeting**: deliver the [right action at the perfect moment](https://docs.criticalmoments.io/conditional-targeting/intro-to-conditions) with over 100 built-in targeting properties. 
  - Example: `device_battery_level > 0.2 && eventCount('app_start') > 3 && app_install_date < now() - duration('24h') && photo_library_permission == 'authorized'`
- **Next Level Privacy**: [100% local, zero data collection](#user-privacy-)

## Table of Contents
1. [How It Works](#how-it-works-)
2. [Feature Overview](#feature-overview-)
3. [Powerful Targeting](#powerful-targeting-)
4. [User Privacy](#user-privacy-)
5. [Quick Start](#quick-start-)
6. [Demo App](#demo-app-)
7. [Documentation](#documentation-)
8. [Contact Us, License & Copyright](#contact-us-)

## How It Works 👩‍💻

### Step 1: Install our SDK

Install our SDK and integrate into your app following our [Quick-Start Guide](https://docs.criticalmoments.io/quick-start). This only takes about 15 minutes.

### Step 2: Create your growth plan in JSON

Yes, really — [a growth plan defined in JSON](https://docs.criticalmoments.io/config-file-structure). You can add features without any additional code, and if you want deeper integrations, custom hooks are available.

Use our [guides](https://docs.criticalmoments.io/guides/reduce-app-churn-with-notifications) to get started with ready-to-deploy and proven growth tactics. We're building a [growing library](https://criticalmoments.io/blog) of examples you can use for inspiration.

### Step 3: Update Anytime, Without App Updates

Once your initial growth plan is deployed, you can update anytime without waiting for app reviews or app store updates. Add growth features without new code, anytime, over the air. Update and tune your user targeting logic, including in past app releases.

## Feature Overview 🔧

- **Notifications**:  Our [smart notifications](https://criticalmoments.io/features/notifications) target delivery to the perfect moment, considering realtime device condition. Start with our templates to [increase activation](https://docs.criticalmoments.io/guides/reduce-app-churn-with-notifications#increase-activation-rate), [reduce churn](https://docs.criticalmoments.io/guides/reduce-app-churn-with-notifications#reduce-long-term-churn), and [custom notifications](https://docs.criticalmoments.io/guides/reduce-app-churn-with-notifications#step-5-add-custom-notification).
- **Improve your App Rating**: use our template to [ask users to rate your app at the perfect moment](https://docs.criticalmoments.io/guides/improve-your-app-store-rating), increasing your rating and rating volume.
- **Optimize Revenue**: ask users to [upgrade at the perfect moment](https://criticalmoments.io/features/grow_revenue), with [over 100 built-in targeting properties](https://docs.criticalmoments.io/conditional-targeting/built-in-properties).
- **In-app Messaging**: add [fully native messaging UI](https://docs.criticalmoments.io/actions-in-app-messaging/actions-overview) without writing any code. Options include [banners](https://docs.criticalmoments.io/actions-in-app-messaging/banners), [modals](https://docs.criticalmoments.io/actions-in-app-messaging/modals), [alerts](https://docs.criticalmoments.io/actions-in-app-messaging/alerts), [browser](https://docs.criticalmoments.io/actions-in-app-messaging/open-link), and more. All [themed](https://docs.criticalmoments.io/themes/theme-overview) to match your brand.
- **Smart Feature Flags**: define [feature flags](https://docs.criticalmoments.io/guides/feature-flags-guide) that can be updated in real-time based on over [100 real-time device conditions \(low battery, has network, permissions, etc\)](https://docs.criticalmoments.io/conditional-targeting/built-in-properties).
- **Over The Air Updates**: Update your growth plan anytime [without app updates](https://docs.criticalmoments.io/remote-control-service). Quickly recover from unexpected bugs, outages, deprecations, and other critical events.


## Powerful Targeting 🎯

With Critical Moments, you can target users with the right actions at the perfect moment:

### Targeting Conditions

Our simple [string-based conditional statements](https://docs.criticalmoments.io/conditional-targeting/intro-to-conditions) can check over [100 built-in properties](https://docs.criticalmoments.io/conditional-targeting/built-in-properties), custom properties, in-app events, and user engagement history. 

Some examples: 
- `eventCount('app_launch') > 5 && latestEventTime('asked_to_subscribe') < now() - duration('72h')`
- `device_model_class == 'iPad' && versionLessThan(app_version, '2.4.1')`
- `camera_permission != 'authorized' && photo_library_permission != 'authorized'`
- `has_watch || location_city == 'Toronto' || has_car_audio || on_call || has_bt_headset || network_connection_type == 'cellular'`
- `weather_condition IN ['Rain', 'Thunderstorms'] || weather_cloud_cover > 0.80`

### Targeting Events

Define exactly when actions should occur, based on in-app event triggers.

The most [commonly needed actions are tracked automatically](https://docs.criticalmoments.io/events/built-in-events). [Add your own custom events](https://docs.criticalmoments.io/events/event-overview) or [properties](https://docs.criticalmoments.io/conditional-targeting/custom-properties) with a single line of code. 

### Local User-Engagement Database 📙

Our SDK automatically starts building an on-device database of user engagement history. Use this database when targeting user-messaging, reviews, notifications, paywalls and more! For example: `eventCount('session_start') > 3 && latestEventTime('asked_to_subscribe') < now() - duration('24h') && !propertyEver('has_paid_subscription', true)`

## User Privacy 🔑🔒

Critical Moments is designed from the ground up for user privacy. All logic is run locally on their own device. We don’t collect any information about your users. A default installation makes zero calls to our servers from the user’s device. Read about our privacy architecture [here](https://criticalmoments.io/blog/how_to_target_users_without_collecting_data).

Note: Some optional services like GeoIP location and weather require external services. These are clearly outlined in our docs. They are completely optional, do not collect user identifiers, and do not store logs long-term.

## Quick Start 🚀

Read our [Quick-Start](https://docs.criticalmoments.io/quick-start) guide to get up and running in minutes. 

## Demo App 

Want to see Critical Moments in action? Download our [demo app from TestFlight](https://testflight.apple.com/join/uSwscwu0) or view the source code [on Github](https://github.com/CriticalMoments/CriticalMoments/tree/main/ios/sample_app).

## Documentation 👩‍💻

Our [documentation](https://docs.criticalmoments.io) includes:

 - [Quick Start Guide](https://docs.criticalmoments.io/quick-start)
 - Guides for [increasing app ratings](https://docs.criticalmoments.io/guides/improve-your-app-store-rating), [reducing churn with notifications](https://docs.criticalmoments.io/guides/reduce-app-churn-with-notifications), and [smart feature flags](https://docs.criticalmoments.io/guides/feature-flags-guide) 
 - [Detailed technical docs](https://docs.criticalmoments.io/concepts-overview)

## Contact Us 👋

We're here to help!

Can't find an answer in our docs? Need help? Find a bug? Feel free to reach out!

- Email: [support@criticalmoments.io](mailto:support@criticalmoments.io) 
- Bug tracker: [GitHub Issues](https://github.com/CriticalMoments/CriticalMoments/issues)

## License ⚖️

Critical Moments requires you to purchase a license to use it in production apps. See our [pricing page](https://criticalmoments.io/pricing) for details. The code of the SDK is fully source-available, and in this repo.

## Copyright ©️

Copyright (c) 2023 Chesterfield Laboratories Inc.

"Critical Moments" and our logos are trademarks. All rights reserved.
