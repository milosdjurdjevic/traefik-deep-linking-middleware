# Deep Linking Traefik Plugin

This project implements a Traefik plugin for mobile redirection. It detects mobile user agents and redirects specific paths to corresponding deep links in a mobile application.

## Overview

The Deep Linking plugin is designed to enhance user experience by seamlessly redirecting mobile users to the appropriate sections of a mobile application based on their requests. It supports various mobile platforms including Android and iOS.

## Installation

To install the Deep Linking plugin, follow the installation instructions available at [Traefik Plugins Catalog](https://plugins.traefik.io/install). You'll need to add the plugin to your Traefik static configuration and enable it in your dynamic configuration.


## Configuration

The plugin requires configuration for the redirects. The configuration is defined in the `Config` struct within the `deep_linking.go` file. You can customize the redirects as needed.

## Usage

Once the plugin is installed and configured, it will automatically redirect mobile users based on their user agent and the requested path.

## Testing

Unit tests for the plugin are located in `deep_linking_test.go`. You can run the tests using:

```
make test
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.