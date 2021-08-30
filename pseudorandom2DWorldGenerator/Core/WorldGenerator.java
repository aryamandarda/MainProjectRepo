package byow.Core;

import byow.TileEngine.TETile;
import byow.TileEngine.Tileset;

import java.io.Serializable;
import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Random;

/**
 * Algorithm behind how the world is generated,
 * rooms are placed, and how hallways connect all rooms.
 *
 * @Source https://gamedevelopment.tutsplus.com/tutorials/
 * how-to-use-bsp-trees-to-generate-game-maps--gamedev-12268
 * @Source https://eskerda.com/bsp-dungeon-generation/
 * @Source https://medium.com/@guribemontero/
 * dungeon-generation-using-binary-space-trees-47d4a668e2d0
 */


public class WorldGenerator implements Serializable {


    private class Node implements Serializable {
        int length;
        int width;
        Position pos;
        Node left, right;
        Room roomContained;
        int orientation = 0; //1- up, down... -1- left, right

        Node(int length, int width, Position pos) {
            this.length = length;
            this.width = width;
            this.pos = pos;
            this.right = null;
            this.left = null;
            this.roomContained = null;
        }

        public boolean isNotEmpty() {
            return this.roomContained != null;
        }
    }

    private static final int MIN_ROOM_SIZE = 4;
    private static final int MAX_NODES = 20;
    private static final int MIN_NODES = 15;
    private static final int MIN_NODE_SIZE = 5;
    private static final int FLOWERS = 15;

    private ArrayList<Room> roomList;       // Stores list of room objects
    private ArrayList<Room> hallwayList;    // Stores list of hallway objects
    private Random randomObject;            // Stores the random object of 'seed'
    private Node root;                      // Stores the root node
    private int maxNoNodes;                 // Stores the max number of nodes
    private ArrayDeque<Node> queue;         // Stores the all the nodes in a queue format
    private TETile[][] tileArray;           // Tile Array representing the World
    private long seed;                      // Seed inputted by the user
    private Position player;                // Position of the player
    private String playerName;              // Player name
    private int points;

    /**
     * Generates the world, fills it with nothing and then
     * initiates our pseudo random world generating algorithm
     */
    public WorldGenerator(long seed, int length, int width) {
        length -= 2;
        this.tileArray = new TETile[width][length];
        this.seed = seed;
        this.playerName = "";
        this.points = 0;
        this.randomObject = new Random(seed);
        maxNoNodes = RandomUtils.uniform(randomObject, MAX_NODES - MIN_NODES + 1) + MIN_NODES;
        this.root = new Node(length, width, new Position(0, 0));
        maxNoNodes -= 1;
        roomList = new ArrayList<Room>();
        hallwayList = new ArrayList<Room>();
        queue = new ArrayDeque<Node>();
        queue.add(root);
        fillWorldWithNothing();
        startAlgorithm(queue, maxNoNodes);
    }

    /*public static void main(String[] args) {
        TERenderer ter = new TERenderer();
        ter.initialize(70, 40);
        long seed = 56367;
        WorldGenerator world = new WorldGenerator(seed, 40, 50);
        ter.renderFrame(world.getWorld());
    }*/

    /**
     * Function that fills tile array with Tilset.NOTHING
     */
    public void fillWorldWithNothing() {
        int height = tileArray[0].length;
        int width = tileArray.length;
        for (int i = 0; i < width; i++) {
            for (int j = 0; j < height; j++) {
                tileArray[i][j] = Tileset.NOTHING;
            }
        }
    }

    /**
     * Our algorithm uses breadth first search to fill a tree based world pseudo randomly
     */
    public void startAlgorithm(ArrayDeque<Node> nodeQueue, int maxNodes) {
        while (maxNodes > 0) {
            Node temp = nodeQueue.remove();
            maxNodes--;
            if (split(temp)) {
                nodeQueue.add(temp.left);
                nodeQueue.add(temp.right);
            }
        }
        /**First generating the rooms and hallways then instantiating them with a for loop*/
        roomsGenerator(root);
        hallwaysGenerator(root);
        //printRooms();
        for (Room room : roomList) {
            drawRoom(room);
        }
        for (Room hallway : hallwayList) {
            drawRoom(hallway);
        }

        /** Creating the player*/
        createPlayer();
        drawTile(player, Tileset.AVATAR);
        addFlowers();
    }

    /*private void printRooms() {
        for (Room room: roomList) {
            System.out.println(room);
        }
    }*/


    /**
     * Splits the tree node into two children nodes.
     * It Randomly decides the split direction
     * and the width / height of the children leaves.
     * The boolean return value indicates whether a node is being split or not
     */

    public boolean split(Node curr) {
        /** Checking if the current node is already split*/
        if (curr.right != null || curr.left != null) {
            return false;
        }
        /** And checking if the size of the room is enough to split*/
        if (curr.length < MIN_NODE_SIZE * 2 && curr.width < MIN_NODE_SIZE * 2) {
            return false;
        }
        /** Then deciding split direction (horizontal or vertical)*/
        int[] values = decideSplitOrientation(curr);
        curr.orientation = values[0];
        int length = values[1];
        /** Splits into two leaves with random size no less than the ROOM_SPLIT_MIN */
        int splitPoint = RandomUtils.uniform(randomObject,
                length - 2 * MIN_NODE_SIZE + 1) + MIN_NODE_SIZE;
        if (curr.orientation == 1) {
            curr.left = new Node(splitPoint, curr.width, curr.pos);
            curr.right = new Node(curr.length - splitPoint, curr.width,
                    new Position(curr.pos.getX(), curr.pos.getY() + splitPoint - 1));
        } else {
            curr.left = new Node(curr.length, splitPoint, curr.pos);
            curr.right = new Node(curr.length, curr.width - splitPoint,
                    new Position(curr.pos.getX() + splitPoint - 1, curr.pos.getY()));
        }
        return true;
    }

    /**
     * Pseudo randomly generating direction to split in with help of
     * seed in the randomObject
     */
    private int[] decideSplitOrientation(Node curr) {
        int[] values = new int[2];
        double decideSplitOrientation = RandomUtils.uniform(randomObject);
        int orientation = 0;
        int length = 0;
        if (decideSplitOrientation < 0.5) {
            orientation = -1; //horizontal
            length = curr.width;
        } else {
            orientation = 1; //vertical
            length = curr.length;
        }
        if (length < MIN_NODE_SIZE * 2) {
            orientation = -orientation;
            if (orientation == 1) {
                length = curr.length;
            } else {
                length = curr.width;
            }
        }
        values[0] = orientation;
        values[1] = length;
        return values;
    }

    /**
     * Recursively creates hallways that connect each two nodes of the same parent.
     * Each hallway starts from the center of a node and goes to the center of the
     * other node.
     */
    public void hallwaysGenerator(Node curr) {
        if (curr.left == null && curr.right == null) {
            return;
        }
        Room hallwayRoom;
        Position newPos;
        if (curr.orientation == -1) {
            int startingPoint = curr.left.pos.getX() + curr.left.width / 2 - 1;
            int endingPoint = curr.right.pos.getX() + curr.right.width / 2 + 1;
            newPos = new Position(startingPoint, curr.pos.getY() + curr.length / 2 - 1);
            hallwayRoom = new Room(newPos, endingPoint - startingPoint + 1, 3);
        } else {
            int startingPoint = curr.left.pos.getY() + curr.left.length / 2 - 1;
            int endingPoint = curr.right.pos.getY() + curr.right.length / 2 + 1;
            newPos = new Position(curr.pos.getX() + curr.width / 2 - 1, startingPoint);
            hallwayRoom = new Room(newPos, 3, endingPoint - startingPoint + 1);
        }
        hallwayList.add(hallwayRoom);
        hallwaysGenerator(curr.left);
        hallwaysGenerator(curr.right);
    }

    /**
     * Recursively creates rooms of the nodes without children.
     */
    public void roomsGenerator(Node curr) {
        if (curr.left == null && curr.right == null) {
            generateRoom(curr);
            if (curr.isNotEmpty()) {
                roomList.add(curr.roomContained);
            }
            return;
        }
        roomsGenerator(curr.left);
        roomsGenerator(curr.right);
    }


    public void generateRoom(Node curr) {
        if (RandomUtils.uniform(randomObject) < 0.30) {
            return;
        }
        /**Finds the biggest possible minimum width and length of a room
         * that can fit in the given Node*/
        int biggestPossibleWidth = Math.max(MIN_ROOM_SIZE, curr.width / 2 + 2);
        int width = RandomUtils.uniform(randomObject,
                curr.width - biggestPossibleWidth + 1) + biggestPossibleWidth;
        int biggestPossibleLength = Math.max(MIN_ROOM_SIZE, curr.length / 2 + 2);
        int length = RandomUtils.uniform(randomObject,
                curr.length - biggestPossibleLength + 1) + biggestPossibleLength;
        Position newPos = curr.pos;
        if (curr.width != width) {
            newPos.shift(RandomUtils.uniform(randomObject,
                    curr.width - width), 0);
        }
        if (curr.length != length) {
            newPos.shift(0, RandomUtils.uniform(randomObject,
                    curr.length - length));
        }
        curr.roomContained = new Room(newPos, width, length);
    }

    /**
     * Draws row by row, an entire room
     */
    public void drawRoom(Room room) {
        for (int y = 0; y < room.getLength(); y++) {
            drawRowOfRoom(room, y);
        }
    }

    /**
     * Draws one row of a given room in the world
     */
    public void drawRowOfRoom(Room room, int row) {
        int y = room.getPos().getY() + row;
        int endXPos = room.getPos().getX() + room.getWidth() - 1;
        if (row == 0 || row == room.getLength() - 1) {
            for (int x = room.getPos().getX(); x <= endXPos; x++) {
                drawTile(new Position(x, y), Tileset.WALL);
            }
            return;
        }
        drawTile(new Position(room.getPos().getX(), y), Tileset.WALL);
        for (int x = room.getPos().getX() + 1; x < endXPos; x++) {
            drawTile(new Position(x, y), Tileset.FLOOR);
        }
        drawTile(new Position(endXPos, y), Tileset.WALL);
    }

    /**
     * Draws a tile at a given position
     */
    private void drawTile(Position p, TETile t) {
        if (!t.equals(Tileset.WALL)
                || !tileArray[p.getX()][p.getY()].equals(Tileset.FLOOR)) {
            tileArray[p.getX()][p.getY()] = t;
        }
    }

    /**
     * Returns the name of the player/ TETile array or sets player name.
     */

    public TETile[][] getWorld() {
        return tileArray;
    }

    public String getPlayerName() {
        return playerName;
    }

    public void setPlayerName(String name) {
        this.playerName = name;
    }

    public int getPoints() {
        return this.points;
    }

    /**
     * Creates player and treasure in different random rooms.
     */
    private void createPlayer() {
        int roomIndex = RandomUtils.uniform(randomObject, roomList.size());
        Room r = roomList.get(roomIndex);
        player = new Position(r.getPos().getX() + r.getWidth() / 2,
                r.getPos().getY() + r.getLength() / 2);
        r.setContainsPlayer();
    }

    public void addFlowers() {
        int flowers = 0;
        int randWidth;
        int randLength;
        while (flowers < FLOWERS) {
            int roomIndex = RandomUtils.uniform(randomObject, roomList.size());
            Room r = roomList.get(roomIndex);
            if (!r.hasPlayer()) {
                int centerX = r.getPos().getX();
                int centerY = r.getPos().getY();
                randWidth = RandomUtils.uniform(randomObject, r.getWidth() - 2) + centerX + 1;
                randLength = RandomUtils.uniform(randomObject, r.getLength() - 2) + centerY + 1;
                tileArray[randWidth][randLength] = Tileset.FLOWER;
                flowers++;
            }
        }
    }

    public void movePlayer(String direction) {
        switch (direction) {
            case "up":
                if (checkCollisions(player, 1, 0, player.getY())) {
                    tileArray[player.getX()][player.getY()] = Tileset.FLOOR;
                    player = player.shift(0, 1);
                    addPoints(player);
                    tileArray[player.getX()][player.getY()] = Tileset.AVATAR;
                }
                break;
            case "down":
                if (checkCollisions(player, -1, 0, player.getY())) {
                    tileArray[player.getX()][player.getY()] = Tileset.FLOOR;
                    player = player.shift(0, -1);
                    addPoints(player);
                    tileArray[player.getX()][player.getY()] = Tileset.AVATAR;
                }
                break;
            case "right":
                if (checkCollisions(player, 1, player.getX(), 0)) {
                    tileArray[player.getX()][player.getY()] = Tileset.FLOOR;
                    player = player.shift(1, 0);
                    addPoints(player);
                    tileArray[player.getX()][player.getY()] = Tileset.AVATAR;
                }
                break;
            case "left":
                if (checkCollisions(player, -1, player.getX(), 0)) {
                    tileArray[player.getX()][player.getY()] = Tileset.FLOOR;
                    player = player.shift(-1, 0);
                    addPoints(player);
                    tileArray[player.getX()][player.getY()] = Tileset.AVATAR;
                }
                break;
            default:
                break;
        }
    }

    private boolean checkCollisions(Position plyr, int direction, int x, int y) {
        if (x == 0) {
            if (tileArray[plyr.getX()][y + direction].equals(Tileset.WALL)) {
                return false;
            }
        } else if (y == 0) {
            if (tileArray[x + direction][plyr.getY()].equals(Tileset.WALL)) {
                return false;
            }
        }
        return true;
    }

    private void addPoints(Position plyr) {
        if (tileArray[plyr.getX()][plyr.getY()].equals(Tileset.FLOWER)) {
            this.points += 5;
        }
    }
}
