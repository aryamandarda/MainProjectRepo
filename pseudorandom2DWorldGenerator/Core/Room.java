package byow.Core;

import java.io.Serializable;

/**Handles creation and manipulation of rooms in the 2d world generated*/

public class Room implements Serializable {
    private Position pos;
    private int width;
    private int length;
    private boolean containsPlayer;

    public Room(Position pos, int width, int length) {
        this.pos = pos;
        this.width = width;
        this.length = length;
        this.containsPlayer  = false;
    }

    public int getWidth() {
        return this.width;
    }

    public int getLength() {
        return this.length;
    }

    public Position getPos() {
        return this.pos;
    }

    @Override
    public String toString() {
        return ("start from " + "(" + pos.getX() + ", " + pos.getY() + ") "
                + "width=" + width + " length=" + length);
    }

    public void setContainsPlayer() {
        this.containsPlayer = true;
    }

    public boolean hasPlayer() {
        return this.containsPlayer;
    }
}
