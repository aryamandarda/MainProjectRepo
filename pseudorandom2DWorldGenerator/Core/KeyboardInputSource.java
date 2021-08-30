package byow.Core;

import byow.TileEngine.TETile;
import byow.TileEngine.Tileset;
import edu.princeton.cs.introcs.StdDraw;

import java.awt.*;
import java.io.Serializable;

import static byow.Core.Engine.*;

/**
 * Allows interacting with the keyboard inputs in accordance with the example in InputDemo
 */

public class KeyboardInputSource implements InputSource, Serializable {

    private int mouseX;
    private int mouseY;
    private TETile[][] worldArray;

    @Override
    public char getNextKey() {
        while (true) {
            if (StdDraw.hasNextKeyTyped()) {
                return Character.toUpperCase(StdDraw.nextKeyTyped());
            }
        }
    }

    public char getNextKey(WorldGenerator world) {
        while (true) {
            StdDraw.setPenColor(Color.gray);
            StdDraw.filledRectangle(25, 29, 25, 1);
            StdDraw.setPenColor(Color.black);
            StdDraw.text(45, 29, world.getPlayerName());
            StdDraw.text(25, 29, "Points: " + world.getPoints());
            mouseX = (int) Math.floor(StdDraw.mouseX()) - 1;
            mouseY = (int) Math.floor(StdDraw.mouseY()) - 1;
            worldArray = world.getWorld();
            if (mouseX < WIDTH - 1 && mouseY < LENGTH - 3) {
                if (mouseX > 0 && mouseY > 0) {
                    TETile temptile = worldArray[mouseX][mouseY];
                    String tileType = tileType(temptile);
                    StdDraw.setPenColor(Color.yellow);
                    StdDraw.setFont(REG_FONT);
                    StdDraw.text(5, 29, tileType);
                    StdDraw.show();
                }
            }
            if (StdDraw.hasNextKeyTyped()) {
                return Character.toUpperCase(StdDraw.nextKeyTyped());
            }
        }
    }

    private String tileType(TETile tile) {
        if (tile.equals(Tileset.WALL)) {
            return "Wall";
        } else if (tile.equals(Tileset.FLOOR)) {
            return "Floor";
        } else if (tile.equals(Tileset.AVATAR)) {
            return "Avatar";
        } else if (tile.equals(Tileset.FLOWER)) {
            return "Flower";
        }
        return "The Great Void";
    }

    @Override
    public boolean possibleNextInput() {
        return true;
    }
}
