package gitlet;

import java.io.File;
import java.io.IOException;
import java.io.Serializable;
import java.util.*;

import static gitlet.Utils.*;

public class CommitTree implements Serializable {

    protected HashMap<String, ArrayList<String>> commitTree;

    public CommitTree() throws IOException {
        commitTree = new HashMap<>();
    }

    public static void treeLog() {
        String commitHash = Utils.readContentsAsString(Repository.HEAD);
        File tempCommitFile = Utils.join(Repository.COMMITS_DIR, commitHash);
        Commit tempCommit = Utils.readObject(tempCommitFile, Commit.class);
        while (tempCommit.parentCommit != null) {
            System.out.println("===");
            System.out.println("commit " + tempCommit.hashName);
            if (tempCommit.parentCommit.size() > 1) {
                System.out.println("Merge: " + tempCommit.parentCommit.get(0).substring(0, 7)
                        + " " + tempCommit.parentCommit.get(1).substring(0, 7));
            }
            tempCommit.printFormattedDate();
            System.out.println(tempCommit.message);
            System.out.println();
            commitHash = tempCommit.parentCommit.get(0);
            tempCommitFile = Utils.join(Repository.COMMITS_DIR, commitHash);
            tempCommit = Utils.readObject(tempCommitFile, Commit.class);
        }
        System.out.println("===");
        System.out.println("commit " + tempCommit.hashName);
        tempCommit.printFormattedDate();
        System.out.println(tempCommit.message);
        System.out.println();
    }

    public static void globalTreeLog() {
        List<String> fileNameList = Utils.plainFilenamesIn(Repository.COMMITS_DIR);
        for (String id: fileNameList) {
            File tempCommitFile = Utils.join(Repository.COMMITS_DIR, id);
            Commit tempCommit = Utils.readObject(tempCommitFile, Commit.class);
            System.out.println("===");
            System.out.println("commit " + tempCommit.hashName);
            tempCommit.printFormattedDate();
            System.out.println(tempCommit.message);
            System.out.println();
        }
    }

    public static void treeFind(String message) {
        int count = 0;
        List<String> fileNameList = Utils.plainFilenamesIn(Repository.COMMITS_DIR);
        for (String id: fileNameList) {
            File tempCommitFile = Utils.join(Repository.COMMITS_DIR, id);
            Commit tempCommit = Utils.readObject(tempCommitFile, Commit.class);
            if (tempCommit.message.equals(message)) {
                System.out.println(id);
                count++;
            }
        }
        if (count == 0) {
            System.out.println("Found no commit with that message.");
            System.exit(0);
        }
    }

    public static void status() {
        System.out.println("=== Branches ===");
        String currentBranch = Utils.readContentsAsString(Repository.CURRENT_BRANCH);
        List<String> branches = Utils.plainFilenamesIn(Repository.BRANCHES_DIR);
        Collections.sort(branches);
        for (String name: branches) {
            if (name.equals(currentBranch)) {
                System.out.print("*");
            }
            System.out.println(name);
        }
        System.out.println();
        System.out.println("=== Staged Files ===");
        StagingArea tempStage = Utils.readObject(Repository.STAGING_AREA, StagingArea.class);
        if (!tempStage.addHashMap.isEmpty()) {
            ArrayList<String> sortedAddStage = new ArrayList<>(tempStage.addHashMap.keySet());
            Collections.sort(sortedAddStage);
            for (String fileName : sortedAddStage) {
                System.out.println(fileName);
            }
        }
        System.out.println();
        System.out.println("=== Removed Files ===");
        if (!tempStage.removeHashMap.isEmpty()) {
            ArrayList<String> sortedRemoveStage = new ArrayList<>(tempStage.removeHashMap.keySet());
            Collections.sort(sortedRemoveStage);
            for (String fileName : sortedRemoveStage) {
                System.out.println(fileName);
            }
        }
        System.out.println();
        File tempCommitFile = join(Repository.COMMITS_DIR, readContentsAsString(Repository.HEAD));
        Commit tempCommit = readObject(tempCommitFile, Commit.class);
        List<String> filesInCWD = plainFilenamesIn(Repository.CWD);
        Collections.sort(filesInCWD);
        System.out.println("=== Modifications Not Staged For Commit ===");
        modificationsHelper(filesInCWD, tempStage, tempCommit);
        System.out.println();
        System.out.println("=== Untracked Files ===");
        for (String fileName: filesInCWD) {
            if (!tempStage.addHashMap.containsKey(fileName)
                    && !tempCommit.blobMap.containsKey(fileName)) {
                System.out.println(fileName);
            }
        }
        System.out.println();
    }

    private static void modificationsHelper(List<String> filesInCWD, StagingArea tempStage,
                                            Commit currentHead) {
        ArrayList<String> filesThatPass = new ArrayList<>();
        for (String fileName: filesInCWD) {
            File tempCWDFile = Utils.join(Repository.CWD, fileName);
            String cwdFileHash = Utils.sha1(readContentsAsString(tempCWDFile));
            if (currentHead.blobMap.containsKey(fileName)
                    && !cwdFileHash.equals(currentHead.blobMap.get(fileName))
                    && !tempStage.addHashMap.containsKey(fileName)) {
                filesThatPass.add(fileName + " (modified)");
            } else if (tempStage.addHashMap.containsKey(fileName)
                    && !cwdFileHash.equals(tempStage.addHashMap.get(fileName))) {
                filesThatPass.add(fileName + " (modified)");
            }
        }
        for (Map.Entry mapElement : tempStage.addHashMap.entrySet()) {
            String fileName = (String) mapElement.getKey();
            File tempCWDFile = Utils.join(Repository.CWD, fileName);
            if (!tempCWDFile.exists()) {
                filesThatPass.add(fileName + " (deleted)");
            }
        }
        for (Map.Entry mapElement : currentHead.blobMap.entrySet()) {
            String fileName = (String) mapElement.getKey();
            File tempCWDFile = Utils.join(Repository.CWD, fileName);
            if (!tempCWDFile.exists() && !tempStage.removeHashMap.containsKey(fileName)) {
                filesThatPass.add(fileName + " (deleted)");
            }
        }
        Collections.sort(filesThatPass);
        for (String fileName: filesThatPass) {
            System.out.println(fileName);
        }
    }
}
