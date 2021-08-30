package gitlet;

import java.io.File;
import java.io.IOException;
import java.time.LocalDateTime;

/** Driver class for Gitlet, a subset of the Git version-control system.
 *  @author Aryaman Darda
 */

public class Main {

    /** Usage: java gitlet.Main ARGS, where ARGS contains
     *  <COMMAND> <OPERAND1> <OPERAND2> ... 
     */
    public static void main(String[] args) throws IOException {
        if(args.length == 0){
            System.out.println("Please enter a command.");
            System.exit(0);
        }
        String firstArg = args[0];
        switch(firstArg) {
            case "init":
                validateArgs("init", args, 1);
                Repository.initDir();
                break;
            case "add":
                initialisedGitletDir();
                validateArgs("add", args, 2);
                Repository.add(args[1]);
                break;
            case "rm":
                initialisedGitletDir();
                validateArgs("rm", args, 2);
                Repository.rm(args[1]);
                break;
            case "commit":
                initialisedGitletDir();
                validateArgs("commit", args, 2);
                LocalDateTime currTime = LocalDateTime.now();
                Repository.commitFiles(args[1], currTime);
                break;
            case "log":
                initialisedGitletDir();
                validateArgs("log", args, 1);
                Repository.log();
                break;
            case "global-log":
                initialisedGitletDir();
                validateArgs("global-log", args, 1);
                Repository.globalLog();
                break;
            case "find":
                initialisedGitletDir();
                validateArgs("find", args, 2);
                Repository.find(args[1]);
                break;
            case "status":
                initialisedGitletDir();
                validateArgs("status", args, 1);
                Repository.statusMethod();
                break;
            case "checkout":
                initialisedGitletDir();
                int method = validateCheckOutArgs("checkout", args);
                if (method == 2) {
                    Repository.checkoutBranch(args[1]);
                }
                else if (method == 3) {
                    Repository.checkoutFile(Utils.readContentsAsString(Repository.HEAD), args[2]);
                }
                else {
                    Repository.checkoutFile(args[1], args[3]);
                }
                break;
            case "branch":
                initialisedGitletDir();
                validateArgs("branch", args, 2);
                Repository.branch(args[1]);
                break;
            case "rm-branch":
                initialisedGitletDir();
                validateArgs("rm-branch", args, 2);
                Repository.rmBranch(args[1]);
                break;
            case "reset":
                initialisedGitletDir();
                validateArgs("reset", args, 2);
                Repository.reset(args[1]);
                break;
            case "merge":
                initialisedGitletDir();
                validateArgs("merge", args, 2);
                Repository.merge(args[1]);
                break;
            case "add-remote":
                initialisedGitletDir();
                validateArgs("add-remote", args, 3);
                Remote.addRemote(args[1], args[2]);
                break;
            case "rm-remote":
                initialisedGitletDir();
                validateArgs("rm-remote", args, 2);
                Remote.rmRemote(args[1]);
                break;
            default:
                System.out.println("No command with that name exists.");
                System.exit(0);
        }
    }

    public static void validateArgs(String cmd, String[] args, int n) {
        if (args.length != n || !(args[0].equals(cmd))) {
            System.out.println("Incorrect operands.");
            System.exit(0);
        }
    }

    public static int validateCheckOutArgs(String cmd, String[] args){
        int arg = 0;
        boolean correct = true;
        if (args.length > 2) {
            if (!args[0].equals(cmd) || (!args[1].equals("--") && !args[2].equals("--"))) {
                correct = false;
            }
        }
        if (args.length == 2 || args.length == 3 || args.length == 4) {
            arg = args.length;
        }
        if (!correct || arg == 0) {
            System.out.println("Incorrect operands.");
            System.exit(0);
        }
        return arg;
    }

    public static void initialisedGitletDir() {
        File workingDir = new File(System.getProperty("user.dir"));
        File GITLET = Utils.join(workingDir, ".gitlet/");
        if (!GITLET.exists()) {
            System.out.println("Not in an initialized Gitlet directory.");
            System.exit(0);
        }
    }
}
