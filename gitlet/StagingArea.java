package gitlet;

import java.io.Serializable;
import java.util.HashMap;

public class StagingArea implements Serializable {

    protected HashMap<String, String> addHashMap;
    protected HashMap<String, String> removeHashMap;

    public StagingArea() {
        addHashMap = new HashMap<>();
        removeHashMap = new HashMap<>();
    }
}
