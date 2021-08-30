package byow.Core;

import java.io.Serializable;

/**
 * Allows interacting with the string argument as inputs in accordance
 * with the example class in InputDemo.*/

public class StringInputSource implements InputSource, Serializable {
    private String input;
    private int index;

    public StringInputSource(String s) {
        index = 0;
        input = s;
    }

    public char getNextKey() {
        char returnChar = input.charAt(index);
        index += 1;
        return returnChar;
    }

    public char getNextKey(WorldGenerator world) {
        char returnChar = input.charAt(index);
        index += 1;
        return returnChar;
    }


    public boolean possibleNextInput() {
        return index < input.length();
    }

    public void forwardString(int endIndex) {
        while (index < endIndex) {
            index++;
        }
    }
}
