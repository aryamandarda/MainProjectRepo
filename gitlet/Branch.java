package gitlet;

import java.io.Serializable;

public class Branch implements Serializable {
    protected String startingPoint;
    protected String name;
    protected String currentPointer;

    public Branch(String startingPoint, String name, String currentPointer) {
        this.startingPoint = startingPoint;
        this.name = name;
        this.currentPointer = currentPointer;
    }
}
