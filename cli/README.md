# VolvoOnCall CLI
A CLI application] written in Go to interact with the Volvo Cars (Volvo On Call) services.

This project was inspired by [molobrakos/volvooncall](https://github.com/molobrakos/volvooncall), and it aims to maintain a certain level of compatibility with it both API and configuration wise.

# Installation
Go to [Releases](https://github.com/theriverman/VolvoOnCall/releases) and download the latest version.

Alternatively, you can install it from source by executing the following command:
```bash
go install github.com/theriverman/VolvoOnCall/cli
```

# Configuration File
Interaction can be simplified by adding your credentials and optionally your default car's VIN to a configuration file at `$HOME/.voc.conf`:
```ini
username: my-volvo-username
password: my-secret-passowrd
defaultCarVin: YV1ABCD00E1234567
# region: your-custom-region
# url: your-custom-api-url
```

Additionally, the used region and url can be modified too. Possible regions are the following:
- "" (e.g. nothing which is the default value)
- na
- cn

# Commands
This section describes the commands available in VolvoOnCall CLI. Each subsection explains a top-level command. See also the results of `voc --help` or just execute `voc` without any commands.

## cars
Lists all cars associated with your account. No additional options.

## lock
Locks the car identified by its VIN.
- `--vin`

Example:
```bash
voc lock --vin YV12ABC3456789
```

## unlock
Unlocks the car identified by its VIN.
- `--vin`

Example:
```bash
voc unlock --vin YV12ABC3456789
```

## heater
Start or stop the heater in the car identified by its VIN.
- `--vin`
- `start`
- `stop`

Example:
```bash
voc heater --vin YV12ABC3456789 start
voc heater --vin YV12ABC3456789 stop
```

## engine
Start or stop the engine in the car identified by its VIN.
- `--vin`
- `start`
- `stop`

Example:
```bash
voc engine --vin YV12ABC3456789 start
voc engine --vin YV12ABC3456789 stop
```

# blink
Flash the turn signals on the car identified by its VIN.
- `--vin`

Example:
```bash
voc blink --vin YV12ABC3456789
```

# honk
Honk the horn on the car identified by its VIN.
- `--vin`

Example:
```bash
voc honk --vin YV12ABC3456789
```

# status
Get a brief overview about a select car.

There are three options to select from:
- get specific attributes using their original JSON value
- get the full original JSON printed to the console
- get a curated overview of the most important parameters (default)

Examples:
```bash
# Returns the most common status parameters only
voc status -vin YV12ABC3456789

# Returns the original JSON
voc status -vin YV12ABC3456789 --json

# Returns only select attributes
voc status -vin YV12ABC3456789 --attributes windows.frontLeftWindowOpen,averageFuelConsumption,averageSpeed
```
For more advanced query options, see the Path Syntax at [https://github.com/tidwall/gjson](https://github.com/tidwall/gjson).

# trips

Examples:
```bash
voc trips -vin YV12ABC3456789
voc trips -vin YV12ABC3456789 --json
```

# register
Save your VolvoOnCall username and password in $HOME/.voc.conf

Example:
```bash
voc register --username "my-volvo-username" --password "my-volvo-password"
```

# Building the Project
The recommended approach to building the project is using [Make](https://en.wikipedia.org/wiki/Make_(software)).

Typical build targets defined in `Makefile` are the following:
  * **build**: builds the project for your system's OS/Architecture. the output file is `./dist/$(BINARY_NAME)$(BINARY_SUFFIX)`
  * **build-darwin**:   builds the project for Darwin/MacOS targeting amd64 and arm64
  * **build-linux**:    builds the project for Linux targeting 386/amd64/arm/arm64
  * **build-windows**:  builds the project for Windows targeting 386/amd64
  * **build-all**:      builds the project for all above declared targets
  * **create-tar**:     creates a tar.gz archive with the contents os `./dist`
  * **clean**:          removes all built binaries and build artefacts

# Contribution
Fork the repository, do your changes, then open a pull request. Thank you!

# Acknowledgements
  * https://github.com/molobrakos/volvooncall
