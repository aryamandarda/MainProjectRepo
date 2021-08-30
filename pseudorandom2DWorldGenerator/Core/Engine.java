package byow.Core;

import byow.TileEngine.TERenderer;
import byow.TileEngine.TETile;
import edu.princeton.cs.introcs.StdDraw;

import java.awt.*;
import java.io.*;

public class Engine implements Serializable {

    static TERenderer ter = new TERenderer();
    /* Feel free to change the width and height. */
    public static final int WIDTH = 50;
    public static final int LENGTH = 30;
    public static final Font HEAD_FONT = new Font("Arial", Font.BOLD, 40);
    public static final Font SUBHEAD_FONT = new Font("Arial", Font.BOLD, 30);
    public static final Font REG_FONT = new Font("Arial", Font.BOLD, 16);

    private WorldGenerator world;
    private String currentState; //options, seeding, play, quit
    private long seed;
    private StringBuilder seedString;
    private StringBuilder inputs;
    private StringBuilder playerName;

    public void restart() {
        world = null;
        currentState = "";
        seed = -9999;
        seedString = new StringBuilder();
        playerName = new StringBuilder();
        inputs = new StringBuilder();
    }

    /**
     * Method used for exploring a fresh world. This method should handle all inputs,
     * including inputs from the main menu.
     */
    public void interactWithKeyboard() {
        ter.initialize(WIDTH, LENGTH, 1, 1);
        while (true) {
            restart();
            drawMenu();
            InputSource inputSource = new KeyboardInputSource();
            currentState = "options";
            while (currentState.equals("options") || currentState.equals("seeding")) {
                parseMenuOptions(inputSource);
            }
            drawWorld();
            while (currentState.equals("play")) {
                parseKeyboardMovement(inputSource);
                drawWorld();
            }
        }
    }

    /*public static void main(String[] args) {
        Engine bruh = new Engine();
        String arg = "ladds";
        ter.initialize(WIDTH, LENGTH, 1, 1);
        TETile[][] worldArray = bruh.interactWithInputString(arg);
        ter.renderFrame(worldArray);
    }*/

    /**
     * Method used for autograding and testing your code. The input string will be a series
     * of characters (for example, "n123sswwdasdassadwas", "n123sss:q", "lwww". The engine should
     * behave exactly as if the user typed these characters into the engine using
     * interactWithKeyboard.
     *
     * Recall that strings ending in ":q" should cause the game to quite save. For example,
     * if we do interactWithInputString("n123sss:q"), we expect the game to run the first
     * 7 commands (n123sss) and then quit and save. If we then do
     * interactWithInputString("l"), we should be back in the exact same state.
     *
     * In other words, both of these calls:
     *   - interactWithInputString("n123sss:q")
     *   - interactWithInputString("lww")
     *
     * should yield the exact same world state as:
     *   - interactWithInputString("n123sssww")
     *
     * @param input the input string to feed to your program
     * @return the 2D TETile[][] representing the state of the world
     */
    public TETile[][] interactWithInputString(String input) {
        // passed in as an argument, and return a 2D tile representation of the
        // world that would have been drawn if the same inputs had been given
        // to interactWithKeyboard().
        //
        // See proj3.byow.InputDemo for a demo of how you can make a nice clean interface
        // that works for many different input types.
        input = input.toUpperCase();
        StringInputSource inputSource;
        if (input.startsWith("L")) {
            inputs = loadWorld(false);
            inputSource = new StringInputSource(inputs.toString() + input);
            seed = Long.parseLong(inputs.substring(1, inputs.indexOf("S")));
            String tempInput = inputs.substring(0, inputs.indexOf("S") + 1);
            inputs = new StringBuilder();
            inputs.append(tempInput);
            inputSource.forwardString(inputs.indexOf("S") + 1);
        } else {
            inputSource = new StringInputSource(input);
            inputs = new StringBuilder();
            inputs.append(input.substring(0, input.indexOf('S') + 1));
            seed = Long.parseLong(input.substring(1, input.indexOf('S')));
            inputSource.forwardString(input.indexOf('S') + 1);
        }
        world = new WorldGenerator(seed, LENGTH, WIDTH);
        if (playerName != null) {
            world.setPlayerName(playerName.toString());
        }
        currentState = "play";
        while (currentState.equals("play") && inputSource.possibleNextInput()) {
            parseKeyboardMovement(inputSource);
        }

        return world.getWorld();
    }

    public void parseKeyboardMovement(InputSource inputSource) {
        char currChar = Character.toUpperCase(inputSource.getNextKey(world));
        switch (currChar) {
            case 'W':
                inputs.append(currChar);
                this.world.movePlayer("up");
                break;
            case 'S':
                inputs.append(currChar);
                this.world.movePlayer("down");
                break;
            case 'D':
                inputs.append(currChar);
                this.world.movePlayer("right");
                break;
            case 'A':
                inputs.append(currChar);
                this.world.movePlayer("left");
                break;
            case ':':
                currChar = Character.toUpperCase(inputSource.getNextKey(world));
                followedByQ(currChar);
                break;
            default:
                break;
        }
    }

    public void parseMenuOptions(InputSource inputSource) {
        char currChar = Character.toUpperCase(inputSource.getNextKey());
        String input;
        switch (currChar) {
            case 'N':
                inputs.append(currChar);
                currentState = "seeding";
                drawNamePrompt("");
                collectName(inputSource);
                drawSeedPrompt("");
                collectSeed(inputSource);
                break;
            case 'L':
                input = loadWorld(true).toString();
                interactWithInputString(input);
                currentState = "play";
                break;
            case 'R':
                currentState = "play";
                StringBuilder inputString = loadWorld(true);
                showReplay(inputString);
                break;
            case 'Q':
                currentState = "quit";
                System.exit(0);
                break;
            default:
                break;
        }
    }

    public void showReplay(StringBuilder input) {
        String tempInput;
        int index = input.indexOf("S");
        while (index < input.length()) {
            tempInput = input.substring(0, index + 1);
            ter.renderFrame(interactWithInputString(tempInput));
            StdDraw.pause(500);
            index++;
        }
        currentState = "quit";
    }

    public void drawMenu() {
        int height = LENGTH / 2;
        int width = WIDTH / 2;
        StdDraw.clear(Color.black);
        StdDraw.setPenColor(Color.white);
        StdDraw.setFont(HEAD_FONT);
        StdDraw.text(width, height + 10, "CS61B: THE GAME");
        StdDraw.setFont(SUBHEAD_FONT);
        StdDraw.text(width, height + 5, "New Game (N)");
        StdDraw.text(width, height + 1, "Load Game (L)");
        StdDraw.text(width, height - 3, "Replay Last Game (R)");
        StdDraw.text(width, height - 7, "Quit (Q)");
        StdDraw.show();
    }

    public void drawNamePrompt(String name) {
        int height = LENGTH / 2;
        int width = WIDTH / 2;
        StdDraw.clear(Color.black);
        StdDraw.setPenColor(Color.white);
        StdDraw.setFont(SUBHEAD_FONT);
        StdDraw.text(width, height + 5, "Please enter a name for the avatar followed by '0': ");
        StdDraw.text(width, height, name);
        StdDraw.show();
    }

    public void drawSeedPrompt(String seedInput) {
        int height = LENGTH / 2;
        int width = WIDTH / 2;
        StdDraw.clear(Color.black);
        StdDraw.setPenColor(Color.white);
        StdDraw.setFont(HEAD_FONT);
        StdDraw.text(width, height + 10, "SEEDING FOR A NEW GAME");
        StdDraw.setFont(SUBHEAD_FONT);
        StdDraw.text(width, height + 5, "Please enter a seed followed by 'S': ");
        StdDraw.text(width, height, seedInput);
        StdDraw.show();
    }

    public void drawWorld() {
        StdDraw.clear();
        ter.renderFrame(world.getWorld());
        StdDraw.setPenColor(Color.gray);
        StdDraw.filledRectangle(25, 29, 25, 1);
        StdDraw.setPenColor(Color.black);
        StdDraw.setFont(REG_FONT);
        StdDraw.text(45, 29, world.getPlayerName());
        StdDraw.text(25, 29, "Points: " + world.getPoints());
        StdDraw.show();
    }

    public void collectSeed(InputSource inputSource) {
        char nextChar = Character.toUpperCase(inputSource.getNextKey());
        inputs.append(nextChar);
        String tempSeedString = "";
        while (nextChar != 'S') {
            if (Character.isDigit(nextChar)) {
                seedString.append(nextChar);
                tempSeedString = seedString.toString();
                drawSeedPrompt(tempSeedString);
            }
            nextChar = Character.toUpperCase(inputSource.getNextKey());
            inputs.append(nextChar);
        }
        currentState = "play";
        seedString.append(nextChar);
        tempSeedString = seedString.toString();
        seed = Long.parseLong(tempSeedString.substring(0, tempSeedString.indexOf('S')));
        world = new WorldGenerator(seed, LENGTH, WIDTH);
        world.setPlayerName(playerName.toString());
        return;
    }

    public void collectName(InputSource inputSource) {
        char nextChar = Character.toUpperCase(inputSource.getNextKey());
        String tempNameString = "";
        while (nextChar != '0') {
            if (!Character.isDigit(nextChar)) {
                playerName.append(nextChar);
                tempNameString = playerName.toString();
                drawNamePrompt(tempNameString);
            }
            nextChar = Character.toUpperCase(inputSource.getNextKey());
        }
    }

    private void followedByQ(char ch) {
        if (ch == 'Q') {
            saveQuit();
        }
    }

    public void saveQuit() {
        File gameFile = new File(System.getProperty("user.dir") + "/savefile.txt");
        File nameFile = new File(System.getProperty("user.dir") + "/namefile.txt");
        try {
            if (!gameFile.exists()) {
                gameFile.createNewFile();
            }
            if (!nameFile.exists()) {
                nameFile.createNewFile();
            }
            FileOutputStream fos1 = new FileOutputStream(gameFile);
            ObjectOutputStream oos1 = new ObjectOutputStream(fos1);
            oos1.writeObject(inputs);
            FileOutputStream fos2 = new FileOutputStream(nameFile);
            ObjectOutputStream oos2 = new ObjectOutputStream(fos2);
            if (playerName != null) {
                oos2.writeObject(playerName.toString());
            }
            currentState = "quit";
        } catch (FileNotFoundException e) {
            System.out.println("file not found");
            System.exit(0);
        } catch (IOException e) {
            e.printStackTrace();
            System.exit(0);
        }
    }

    public StringBuilder loadWorld(boolean nameExists) {
        FileInputStream fis;
        ObjectInputStream ois;
        File gameFile = new File(System.getProperty("user.dir") + "/savefile.txt");
        File nameFile = new File(System.getProperty("user.dir") + "/namefile.txt");
        if (gameFile.exists()) {
            try {
                fis = new FileInputStream(gameFile);
                ois = new ObjectInputStream(fis);
                inputs = (StringBuilder) ois.readObject();
            } catch (ClassNotFoundException | IOException e) {
                e.printStackTrace();
                System.exit(0);
            }
        }
        if (nameFile.exists()) {
            if (nameExists) {
                try {
                    fis = new FileInputStream(nameFile);
                    ois = new ObjectInputStream(fis);
                    String name = (String) ois.readObject();
                    playerName.append(name);
                } catch (ClassNotFoundException | IOException e) {
                    e.printStackTrace();
                    System.exit(0);
                }
            }
        } else {
            System.exit(0);
        }
        return inputs;
    }
}
