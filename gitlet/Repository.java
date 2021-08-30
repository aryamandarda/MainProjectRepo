package gitlet;

import java.io.File;
import java.io.IOException;
import java.io.Serializable;
import java.time.LocalDateTime;
import java.util.*;

import static gitlet.Utils.*;

/** Represents a gitlet repository.
 *  Main class of the program. Contains implementations for most of the command in GITLET.
 *  Interacts directly with the .gitlet repository and the CWD.
 *  @author Aryaman Darda
 */

public class Repository implements Serializable {
    /**
     * List all instance variables of the Repository class here with a useful
     * comment above them describing what that variable represents and how that
     * variable is used. We've provided two examples for you.
     */

    /** The current working directory. */
    public static final File CWD = new File(System.getProperty("user.dir"));
    /** The .gitlet directory. */
    public static final File GITLET_DIR = join(CWD, ".gitlet/");
    public static final File STAGING_AREA = join(GITLET_DIR, ".stagingArea");
    protected static File HEAD = join(GITLET_DIR, ".HEAD");
    static final File BLOBS_DIR = Utils.join(Repository.GITLET_DIR, ".blobs/");
    public static final File COMMITS_DIR = Utils.join(Repository.GITLET_DIR, ".commits/");
    public static final File COMMIT_TREE = Utils.join(Repository.GITLET_DIR, ".commitTree");
    public static final File BRANCHES_DIR = Utils.join(Repository.GITLET_DIR, ".branches/");
    public static final File CURRENT_BRANCH = Utils.join(Repository.GITLET_DIR, ".currentBranch");
    public static final File REMOTE = join(Repository.GITLET_DIR, ".remote");

    public static void initDir() throws IOException {
        if (GITLET_DIR.exists()) {
            System.out.println("A Gitlet version-control system already "
                    + "exists in the current directory.");
            System.exit(0);
        }
        GITLET_DIR.mkdir();
        STAGING_AREA.createNewFile();
        BLOBS_DIR.mkdir();
        BRANCHES_DIR.mkdir();
        REMOTE.createNewFile();
        StagingArea stagingArea = new StagingArea();
        Utils.writeObject(STAGING_AREA, stagingArea);
        COMMITS_DIR.mkdir();
        Commit initialCommit = new Commit();
        CommitTree commitT = new CommitTree();
        commitT.commitTree.put(initialCommit.getHashName(), null);
        File commitFile = join(COMMITS_DIR, initialCommit.getHashName());
        commitFile.createNewFile();
        COMMIT_TREE.createNewFile();
        Utils.writeObject(COMMIT_TREE, commitT);
        HEAD.createNewFile();
        Branch master = new Branch(initialCommit.getHashName(), "master",
                initialCommit.getHashName());
        File masterFile = join(BRANCHES_DIR, master.name);
        masterFile.createNewFile();
        CURRENT_BRANCH.createNewFile();
        Utils.writeContents(CURRENT_BRANCH, master.name);
        Utils.writeContents(HEAD, initialCommit.getHashName());
        Utils.writeObject(masterFile, master);
        Utils.writeObject(commitFile, initialCommit);
    }

    public static void add(String filename) throws IOException {
        File tempFile = join(CWD, filename);
        if (!tempFile.exists()) {
            System.out.println("File does not exist.");
            System.exit(0);
        }
        String fileContents = readContentsAsString(tempFile);
        String hashFile = sha1(fileContents);
        StagingArea tempStage = Utils.readObject(STAGING_AREA, StagingArea.class);
        File tempHeadCommit = join(COMMITS_DIR, readContentsAsString(HEAD));
        Commit headCommitCopy = readObject(tempHeadCommit, Commit.class);
        if (tempStage.addHashMap.containsKey(filename)) {
            if (headCommitCopy.blobMap.containsKey(filename)) {
                if (hashFile.equals(headCommitCopy.blobMap.get(filename))) {
                    tempStage.addHashMap.remove(filename);
                }
            } else {
                tempStage.addHashMap.replace(filename, hashFile);
                File blobFile = join(BLOBS_DIR, hashFile);
                blobFile.createNewFile();
                String contents = readContentsAsString(tempFile);
                writeContents(blobFile, contents);
            }
        } else {
            if (!hashFile.equals(headCommitCopy.blobMap.get(filename))) {
                tempStage.addHashMap.put(filename, hashFile);
                File blobFile = join(BLOBS_DIR, hashFile);
                blobFile.createNewFile();
                String contents = readContentsAsString(tempFile);
                writeContents(blobFile, contents);
            }
        }
        if (tempStage.removeHashMap.containsKey(filename)) {
            tempStage.removeHashMap.remove(filename);
        }
        Utils.writeObject(STAGING_AREA, tempStage);
    }

    public static void rm(String filename) {
        StagingArea tempStage = Utils.readObject(STAGING_AREA, StagingArea.class);
        File tempHeadCommit = join(COMMITS_DIR, readContentsAsString(HEAD));
        Commit headCommitCopy = readObject(tempHeadCommit, Commit.class);
        if (!headCommitCopy.blobMap.containsKey(filename)
                && !tempStage.addHashMap.containsKey(filename)) {
            System.out.println("No reason to remove the file.");
            System.exit(0);
        }
        if (headCommitCopy.blobMap.containsKey(filename)) {
            String hashFile = headCommitCopy.blobMap.get(filename);
            tempStage.removeHashMap.put(filename, hashFile);
            File removal = join(CWD, filename);
            if (removal.exists()) {
                Utils.restrictedDelete(removal);
            }
        }
        if (tempStage.addHashMap.containsKey(filename)) {
            tempStage.addHashMap.remove(filename);
        }
        Utils.writeObject(STAGING_AREA, tempStage);
    }

    public static void commitFiles(String message, LocalDateTime dateTime) throws IOException {
        if (message.length() == 0) {
            System.out.println("Please enter a commit message.");
            System.exit(0);
        }
        StagingArea tempStage = readObject(STAGING_AREA, StagingArea.class);
        if (tempStage.addHashMap.isEmpty() && tempStage.removeHashMap.isEmpty()) {
            System.out.println("No changes added to the commit.");
            System.exit(0);
        }
        Commit commitObject = new Commit(message, dateTime);
        File tempParentCommitFile = join(COMMITS_DIR, commitObject.parentCommit.get(0));
        Commit parentCommitTemp = readObject(tempParentCommitFile, Commit.class);
        if (!parentCommitTemp.blobMap.isEmpty()) {
            commitObject.blobMap.putAll(parentCommitTemp.blobMap);
            HashMap<String, String> tempMap = new HashMap<>(commitObject.blobMap);
            for (Map.Entry mapElement : tempMap.entrySet()) {
                String filename = (String) mapElement.getKey();
                if (tempStage.removeHashMap.containsKey(filename)) {
                    commitObject.blobMap.remove(filename);
                } else if (tempStage.addHashMap.containsKey(filename)) {
                    commitObject.blobMap.replace(filename, tempStage.addHashMap.get(filename));
                    tempStage.addHashMap.remove(filename);
                }
            }
        }
        commitObject.blobMap.putAll(tempStage.addHashMap);
        tempStage.addHashMap.clear();
        tempStage.removeHashMap.clear();
        commitObject.setHashName();
        CommitTree tempTree = readObject(COMMIT_TREE, CommitTree.class);
        tempTree.commitTree.put(commitObject.getHashName(), commitObject.parentCommit);
        writeObject(COMMIT_TREE, tempTree);
        writeContents(HEAD, commitObject.getHashName());
        File tempCurrentBranch = join(BRANCHES_DIR, readContentsAsString(CURRENT_BRANCH));
        Branch tempCurrent = readObject(tempCurrentBranch, Branch.class);
        tempCurrent.currentPointer = commitObject.getHashName();
        writeObject(tempCurrentBranch, tempCurrent);
        writeObject(STAGING_AREA, tempStage);
        File commitFile = join(COMMITS_DIR, commitObject.getHashName());
        commitFile.createNewFile();
        writeObject(commitFile, commitObject);
    }

    public static void log() {
        CommitTree.treeLog();
    }

    public static void globalLog() {
        CommitTree.globalTreeLog();
    }

    public static void find(String message) {
        CommitTree.treeFind(message);
    }

    public static void statusMethod() {
        CommitTree.status();
    }

    public static void checkoutFile(String commitID, String filename) throws IOException {
        String commitHash = commitIDExists(commitID);
        File tempCommitFile = Utils.join(COMMITS_DIR, commitHash);
        Commit tempCommit = Utils.readObject(tempCommitFile, Commit.class);
        if (!tempCommit.blobMap.containsKey(filename)) {
            System.out.println("File does not exist in that commit.");
            System.exit(0);
        }
        String fileHash = tempCommit.blobMap.get(filename);
        File tempFile = join(BLOBS_DIR, fileHash);
        String tempFileContents = readContentsAsString(tempFile);
        File currentWorkingFile = join(CWD, filename);
        if (!currentWorkingFile.exists()) {
            currentWorkingFile.createNewFile();
        }
        writeContents(currentWorkingFile, tempFileContents);
    }

    private static String commitIDExists(String commitID) {
        boolean exists = false;
        String commitHash = "";
        CommitTree tempTree = readObject(COMMIT_TREE, CommitTree.class);
        for (Map.Entry mapElement : tempTree.commitTree.entrySet()) {
            String id = (String) mapElement.getKey();
            String shortID = id.substring(0, 6);
            if (commitID.contains(shortID) || commitID.equals(id)) {
                commitHash = id;
                exists = true;
            }
        }
        if (!exists) {
            System.out.println("No commit with that id exists.");
            System.exit(0);
        }
        return commitHash;
    }

    private static void noUntrackedFiles(List<String> filesInCWD,
                                         StagingArea tempStage, Commit tempCommit) {
        for (String fileName: filesInCWD) {
            if (!tempStage.addHashMap.containsKey(fileName)
                    && !tempCommit.blobMap.containsKey(fileName)) {
                System.out.println("There is an untracked file in the way; "
                        + "delete it, or add and commit it first.");
                System.exit(0);
            }
        }
    }

    public static void checkoutBranch(String branchName) throws IOException {
        File tempBranchFile = join(BRANCHES_DIR, branchName);
        if (!tempBranchFile.exists()) {
            System.out.println("No such branch exists.");
            System.exit(0);
        }
        Branch tempBranch = readObject(tempBranchFile, Branch.class);
        String currentBranch = readContentsAsString(CURRENT_BRANCH);
        if (currentBranch.equals(branchName)) {
            System.out.println("No need to checkout the current branch.");
            System.exit(0);
        }
        File tempCommitFile = join(COMMITS_DIR, readContentsAsString(HEAD));
        Commit tempCommit = readObject(tempCommitFile, Commit.class);
        StagingArea tempStage = readObject(STAGING_AREA, StagingArea.class);
        List<String> filesInCWD = plainFilenamesIn(CWD);
        noUntrackedFiles(filesInCWD, tempStage, tempCommit);
        File tempBranchCommitFile = join(COMMITS_DIR, tempBranch.currentPointer);
        Commit tempBranchHeadCommit = readObject(tempBranchCommitFile, Commit.class);
        for (Map.Entry mapElement : tempBranchHeadCommit.blobMap.entrySet()) {
            String fileName = (String) mapElement.getKey();
            String fileID = tempBranchHeadCommit.blobMap.get(fileName);
            File blobFile = join(BLOBS_DIR, fileID);
            String contents = readContentsAsString(blobFile);
            File cwdFile = join(CWD, fileName);
            if (!cwdFile.exists()) {
                cwdFile.createNewFile();
            }
            writeContents(cwdFile, contents);
        }
        for (String fileName: filesInCWD) {
            File cwdFile = join(CWD, fileName);
            if (!tempBranchHeadCommit.blobMap.containsKey(fileName)) {
                restrictedDelete(cwdFile);
            }
        }
        tempStage.addHashMap.clear();
        tempStage.removeHashMap.clear();
        writeContents(HEAD, tempBranch.currentPointer);
        writeContents(CURRENT_BRANCH, tempBranch.name);
        writeObject(STAGING_AREA, tempStage);
    }

    public static void branch(String branchName) throws IOException {
        File newBranchFile = join(BRANCHES_DIR, branchName);
        if (newBranchFile.exists()) {
            System.out.println("A branch with that name already exists.");
            System.exit(0);
        }
        newBranchFile.createNewFile();
        Branch newBranch = new Branch(readContentsAsString(HEAD), branchName,
                readContentsAsString(HEAD));
        writeObject(newBranchFile, newBranch);
    }

    public static void rmBranch(String branchName) {
        File tempBranchFile = join(BRANCHES_DIR, branchName);
        if (!tempBranchFile.exists()) {
            System.out.println("A branch with that name does not exist.");
            System.exit(0);
        }
        String currentBranch = readContentsAsString(CURRENT_BRANCH);
        if (currentBranch.equals(branchName)) {
            System.out.println("Cannot remove the current branch.");
            System.exit(0);
        }
        tempBranchFile.delete();
    }

    public static void reset(String commitID) throws IOException {
        String commitHash = commitIDExists(commitID);
        File tempCommitFile = join(COMMITS_DIR, commitHash);
        Commit tempCommit = Utils.readObject(tempCommitFile, Commit.class);
        File tempHeadCommitFile = join(COMMITS_DIR, readContentsAsString(HEAD));
        Commit tempHeadCommit = readObject(tempHeadCommitFile, Commit.class);
        StagingArea tempStage = readObject(STAGING_AREA, StagingArea.class);
        List<String> filesInCWD = plainFilenamesIn(CWD);
        noUntrackedFiles(filesInCWD, tempStage, tempHeadCommit);
        for (Map.Entry mapElement : tempCommit.blobMap.entrySet()) {
            String fileName = (String) mapElement.getKey();
            checkoutFile(commitHash, fileName);
        }
        for (String fileName: filesInCWD) {
            File cwdFile = join(CWD, fileName);
            if (!tempCommit.blobMap.containsKey(fileName)) {
                restrictedDelete(cwdFile);
            }
        }
        File tempBranchFile = join(BRANCHES_DIR, readContentsAsString(CURRENT_BRANCH));
        Branch tempBranch = readObject(tempBranchFile, Branch.class);
        tempBranch.currentPointer = commitHash;
        tempStage.addHashMap.clear();
        tempStage.removeHashMap.clear();
        writeContents(HEAD, commitHash);
        writeObject(tempBranchFile, tempBranch);
        writeObject(STAGING_AREA, tempStage);
    }

    public static void merge(String branchName) throws IOException {
        StagingArea tempStage = readObject(STAGING_AREA, StagingArea.class);
        if (!tempStage.addHashMap.isEmpty() || !tempStage.removeHashMap.isEmpty()) {
            System.out.println("You have uncommitted changes.");
            System.exit(0);
        }
        File tempBranchFile = join(BRANCHES_DIR, branchName);
        if (!tempBranchFile.exists()) {
            System.out.println("A branch with that name does not exist.");
            System.exit(0);
        }
        Branch tempBranch = readObject(tempBranchFile, Branch.class);
        String currentBranch = readContentsAsString(CURRENT_BRANCH);
        if (currentBranch.equals(branchName)) {
            System.out.println("Cannot merge a branch with itself.");
            System.exit(0);
        }
        File tempHeadCommitFile = join(COMMITS_DIR, readContentsAsString(HEAD));
        Commit tempHeadCommit = readObject(tempHeadCommitFile, Commit.class);
        List<String> filesInCWD = plainFilenamesIn(CWD);
        noUntrackedFiles(filesInCWD, tempStage, tempHeadCommit);
        File givenBranchHeadCommitFile = join(COMMITS_DIR, tempBranch.currentPointer);
        Commit givenBranchHeadCommit = readObject(givenBranchHeadCommitFile, Commit.class);
        Commit splitPointCommit = findSplitPoint(givenBranchHeadCommit, tempHeadCommit);
        File tempCurrentBranchFile = join(BRANCHES_DIR, currentBranch);
        Branch tempCurrentBranch = readObject(tempCurrentBranchFile, Branch.class);
        mergeSpecialCases(splitPointCommit, tempBranch, tempCurrentBranch);
        splitPointFileCheck(tempHeadCommit, givenBranchHeadCommit, splitPointCommit);
        givenBranchFileCheck(tempHeadCommit, givenBranchHeadCommit, splitPointCommit);
        commitFiles("Merged " + branchName
                + " into " + currentBranch + ".", LocalDateTime.now());
        tempHeadCommitFile = join(COMMITS_DIR, readContentsAsString(HEAD));
        tempHeadCommit = readObject(tempHeadCommitFile, Commit.class);
        tempHeadCommit.parentCommit.add(tempBranch.currentPointer);
        CommitTree tempCommitTree = readObject(COMMIT_TREE, CommitTree.class);
        ArrayList<String> tempArray =
                tempCommitTree.commitTree.get(readContentsAsString(HEAD));
        tempArray.add(tempBranch.currentPointer);
        tempCommitTree.commitTree
                .replace(readContentsAsString(HEAD), tempArray);
        writeObject(tempHeadCommitFile, tempHeadCommit);
        writeObject(COMMIT_TREE, tempCommitTree);
    }

    private static void splitPointFileCheck(Commit tempHeadCommit,
                                            Commit givenBranchHeadCommit,
                                            Commit splitPointCommit) throws IOException {
        for (Map.Entry mapElement : splitPointCommit.blobMap.entrySet()) {
            String fileName = (String) mapElement.getKey();
            String fileHash = (String) mapElement.getValue();
            if (tempHeadCommit.blobMap.containsKey(fileName)
                    && givenBranchHeadCommit.blobMap.containsKey(fileName)) {
                if (fileHash.equals(tempHeadCommit.blobMap.get(fileName))
                        && !fileHash.equals(givenBranchHeadCommit.blobMap.get(fileName))) {
                    checkoutFile(givenBranchHeadCommit.hashName, fileName);
                    Repository.add(fileName);
                } else if (fileHash.equals(givenBranchHeadCommit.blobMap.get(fileName))
                        && !fileHash.equals(tempHeadCommit.blobMap.get(fileName))) {
                    continue;
                } else if (!givenBranchHeadCommit.blobMap.get(fileName)
                        .equals(tempHeadCommit.blobMap.get(fileName))) {
                    mergeConflict1(givenBranchHeadCommit, tempHeadCommit, fileName);
                }
            }
            if (fileHash.equals(tempHeadCommit.blobMap.get(fileName))
                    && !givenBranchHeadCommit.blobMap.containsKey(fileName)) {
                rm(fileName);
            }
            if (!tempHeadCommit.blobMap.containsKey(fileName)
                    && !givenBranchHeadCommit.blobMap.containsKey(fileName)) {
                File tempFile = join(CWD, fileName);
                if (tempFile.exists()) {
                    continue;
                }
            } else if (!tempHeadCommit.blobMap.containsKey(fileName)
                    && !fileHash.equals(givenBranchHeadCommit.blobMap.get(fileName))) {
                mergeConflict3(givenBranchHeadCommit, fileName);
            } else if (!fileHash.equals(tempHeadCommit.blobMap.get(fileName))
                    && !givenBranchHeadCommit.blobMap.containsKey(fileName)) {
                mergeConflict2(tempHeadCommit, fileName);
            }
        }
    }

    private static void givenBranchFileCheck(Commit tempHeadCommit,
                                             Commit givenBranchHeadCommit,
                                             Commit splitPointCommit) throws IOException {
        for (Map.Entry mapElement : givenBranchHeadCommit.blobMap.entrySet()) {
            String fileName = (String) mapElement.getKey();
            if (!splitPointCommit.blobMap.containsKey(fileName)
                    && !tempHeadCommit.blobMap.containsKey(fileName)) {
                checkoutFile(givenBranchHeadCommit.hashName, fileName);
                add(fileName);
            } else if (!splitPointCommit.blobMap.containsKey(fileName)
                    && tempHeadCommit.blobMap.containsKey(fileName)) {
                if (!givenBranchHeadCommit.blobMap.get(fileName)
                        .equals(tempHeadCommit.blobMap.get(fileName))) {
                    mergeConflict1(givenBranchHeadCommit, tempHeadCommit, fileName);
                }
            }
        }
    }
    
    private static void mergeSpecialCases(Commit splitPointCommit,
                                          Branch tempBranch,
                                          Branch tempCurrentBranch) throws IOException {
        if (tempBranch.currentPointer.equals(splitPointCommit.hashName)) {
            System.out.println("Given branch is an ancestor of the current branch.");
            System.exit(0);
        }
        if (tempCurrentBranch.currentPointer.equals(splitPointCommit.hashName)) {
            checkoutBranch(tempBranch.name);
            System.out.println("Current branch fast-forwarded.");
            System.exit(0);
        }
    }

    private static void mergeConflict1(Commit givenBranchHeadCommit,
                                      Commit tempHeadCommit,
                                      String fileName) throws IOException {
        System.out.println("Encountered a merge conflict.");
        File givenBranchFile = join(BLOBS_DIR,
                givenBranchHeadCommit.blobMap.get(fileName));
        File currentBranchFile = join(BLOBS_DIR,
                tempHeadCommit.blobMap.get(fileName));
        String givenBranchFileContents =
                readContentsAsString(givenBranchFile);
        String currentBranchFileContents =
                readContentsAsString(currentBranchFile);
        File modifiedFile = join(CWD, fileName);
        writeContents(modifiedFile, "<<<<<<< HEAD\n"
                + currentBranchFileContents + "=======\n"
                + givenBranchFileContents + ">>>>>>>\n");
        Repository.add(fileName);
    }

    private static void mergeConflict2(Commit tempHeadCommit, String fileName) throws IOException {
        System.out.println("Encountered a merge conflict.");
        File currentBranchFile = join(BLOBS_DIR,
                tempHeadCommit.blobMap.get(fileName));
        String currentBranchFileContents =
                readContentsAsString(currentBranchFile);
        File modifiedFile = join(CWD, fileName);
        writeContents(modifiedFile, "<<<<<<< HEAD\n"
                + currentBranchFileContents + "=======\n"
                + ">>>>>>>\n");
        Repository.add(fileName);
    }

    private static void mergeConflict3(Commit tempHeadCommit, String fileName) throws IOException {
        System.out.println("Encountered a merge conflict.");
        File givenBranchFile = join(BLOBS_DIR,
                tempHeadCommit.blobMap.get(fileName));
        String givenBranchFileContents =
                readContentsAsString(givenBranchFile);
        File modifiedFile = join(CWD, fileName);
        writeContents(modifiedFile, "<<<<<<< HEAD\n"
                + "=======\n" + givenBranchFileContents
                + ">>>>>>>\n");
        Repository.add(fileName);
    }

    private static Commit findSplitPoint(Commit givenBranchHeadCommit,
                                        Commit currentBranchHeadCommit) {
        ArrayList<Commit> givenAncestors = findAncestors(givenBranchHeadCommit);
        ArrayList<Commit> currentAncestors = findAncestors(currentBranchHeadCommit);
        for (Commit c1: currentAncestors) {
            for (Commit c2: givenAncestors) {
                if (c1.hashName.equals(c2.hashName)) {
                    return c1;
                }
            }
        }
        return null;
    }

    private static ArrayList<Commit> findAncestors(Commit commit) {
        ArrayDeque<Commit> fringe = new ArrayDeque<>();
        ArrayList<Commit> ancestors = new ArrayList<>();
        fringe.add(commit);
        while (!fringe.isEmpty()) {
            Commit tempCommit = fringe.poll();
            ancestors.add(tempCommit);
            if (tempCommit.parentCommit != null) {
                for (String hash : tempCommit.parentCommit) {
                    File tempAncestorCommitFile = join(COMMITS_DIR, hash);
                    Commit tempAncestorCommit = readObject(tempAncestorCommitFile, Commit.class);
                    fringe.add(tempAncestorCommit);
                }
            }
        }
        return ancestors;
    }
}
