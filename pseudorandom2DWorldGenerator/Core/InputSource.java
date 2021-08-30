package byow.Core;

/**
 * Contains the two methods that allow interacting with the inputs in accordance with
 * the example interface in InputDemo.
 */

public interface InputSource {
    char getNextKey();
    char getNextKey(WorldGenerator world);
    boolean possibleNextInput();
}
