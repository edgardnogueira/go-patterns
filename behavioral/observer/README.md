# Observer Pattern

## Intent

The Observer pattern defines a one-to-many dependency between objects so that when one object changes state, all its dependents are notified and updated automatically.

## Explanation

This implementation demonstrates a weather monitoring system where a weather station (subject) tracks temperature, humidity, and pressure, and various displays (observers) subscribe to receive updates when these measurements change.

## Structure

- **Subject**: Interface that defines methods for attaching, detaching, and notifying observers (WeatherSubject)
- **ConcreteSubject**: Implements the Subject interface and maintains state that Observers are interested in (WeatherStation)
- **Observer**: Interface that defines an update method for objects that should be notified of changes in a Subject (WeatherObserver)
- **ConcreteObserver**: Implements the Observer interface to keep its state consistent with the Subject's state (CurrentConditionsDisplay, StatisticsDisplay, ForecastDisplay)

## When to Use

- When a change to one object requires changing others, and you don't know how many objects need to be changed
- When an object should be able to notify other objects without making assumptions about who these objects are
- When an abstraction has two aspects, one dependent on the other, and encapsulating these aspects in separate objects lets you vary and reuse them independently

## Benefits

- Provides a loosely coupled design between objects that interact
- Enables broadcast communication
- Supports the principle of open/closed design (open for extension, closed for modification)
- Simplifies maintenance by centralizing the update logic in observable objects
