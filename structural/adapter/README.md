# Adapter Pattern

## Intent

The Adapter pattern converts the interface of a class into another interface clients expect. It allows classes to work together that couldn't otherwise because of incompatible interfaces.

## Explanation

This implementation demonstrates a scenario where we have a media player application that can only play MP3 files, but we want to extend it to play other formats like FLAC and WAV without modifying the existing code. We use adapters to convert the interfaces of different media formats into the interface the media player expects.

## Structure

- **Target**: The interface that clients use (MediaPlayer)
- **Adaptee**: The interface that needs adapting (AdvancedMediaPlayer)
- **Adapter**: The class that adapts the Adaptee to the Target (MediaAdapter)
- **Client**: The class that interacts with the Target (AudioPlayer)

## When to Use

- When you want to use an existing class, but its interface doesn't match the one you need
- When you want to create a reusable class that cooperates with classes that don't necessarily have compatible interfaces
- When you need to use several existing subclasses, but it's impractical to adapt their interfaces by subclassing every one

## Benefits

- Allows two incompatible interfaces to work together
- Improves reusability of older code
- Increases flexibility by decoupling client from implementation
- Enables client code to work with unforeseen classes
