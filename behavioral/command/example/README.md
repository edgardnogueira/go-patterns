# Command Pattern Example: Smart Home Automation System

This example demonstrates a practical implementation of the Command pattern in a smart home automation system. The application allows you to control various home devices through a unified interface.

## Overview

The Smart Home Automation System demonstrates:

1. The use of various commands to control different devices
2. Using a remote control as an invoker for commands
3. Creating and executing scene macros (composite commands)
4. Scheduling commands for future execution
5. Command history and undo functionality

## Devices

The system includes several smart home devices:
- Lights (with on/off and dimming capabilities)
- Thermostat (with temperature and mode settings)
- Audio System (with on/off and playback capabilities)
- Garage Door (with open/close and light functionality)

## Commands

Various commands are implemented to control these devices:
- Light commands (on, off, dim)
- Thermostat commands (set temperature, set mode)
- Audio system commands (on, off, play, stop)
- Garage door commands (open, close, light on/off)

## Features

1. **Remote Control**: Simulates a physical remote with multiple buttons
2. **Scenes**: Predefined sets of commands executed together (morning, evening, away)
3. **Scheduling**: Queue commands to be executed at a future time
4. **Undo**: Revert the last executed command
5. **Status Monitoring**: Check the status of all devices

## Running the Example

From the root of the repository, run:

```bash
cd behavioral/command/example
go run main.go
```

## Using the Application

The application provides a text-based menu with the following options:

1. **Show Device Status**: Display the current state of all devices
2. **Use Remote Control**: Simulate pressing buttons on a remote control
3. **Activate Scene**: Execute a predefined set of commands (Morning, Evening, Away)
4. **Schedule Command**: Add a command to the execution queue with a delay
5. **Execute Scheduled Commands**: Run any commands that are due for execution
6. **Undo Last Command**: Revert the most recently executed command
7. **Exit**: Quit the application

## Example Workflow

1. Check the initial status of all devices
2. Use the remote control to turn on some lights
3. Activate the "Evening Scene" to set up the house for evening relaxation
4. Schedule the lights to turn off in 30 seconds
5. After some time, execute the scheduled commands
6. Use the undo feature to revert the last command

This example shows how the Command pattern makes it easy to:
- Decouple the invoker (remote) from the receivers (devices)
- Create complex operations from simple ones (scenes)
- Support scheduling and undo/redo functionality
- Add new commands without changing existing code
