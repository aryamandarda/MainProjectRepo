package gitlet;

import java.io.File;
import java.io.Serializable;
import java.util.HashMap;

import static gitlet.Utils.*;

public class Remote implements Serializable {
    protected HashMap<String, String> remotes = new HashMap<>();

    public static void addRemote(String remoteName, String path) {
        path.replaceAll("/", File.separator);
        File tempRemoteFile = join(Repository.GITLET_DIR, ".remote");
        Remote tempRemote = readObject(tempRemoteFile, Remote.class);
        if (tempRemote.remotes.containsKey(remoteName)) {
            System.out.println("A remote with that name already exists.");
            System.exit(0);
        }
        tempRemote.remotes.put(remoteName, path);
        Utils.writeObject(tempRemoteFile, tempRemote);
    }

    public static void rmRemote(String remoteName) {
        File tempRemoteFile = join(Repository.GITLET_DIR, ".remote");
        Remote tempRemote = readObject(tempRemoteFile, Remote.class);
        if (!tempRemote.remotes.containsKey(remoteName)) {
            System.out.println("A remote with that name does not exist.");
            System.exit(0);
        }
        tempRemote.remotes.remove(remoteName);
        Utils.writeObject(tempRemoteFile, tempRemote);
    }

    public static void push(String remoteName, String remoteBranchName) {
        File tempRemoteFile = join(Repository.GITLET_DIR, ".remote");
        Remote tempRemote = readObject(tempRemoteFile, Remote.class);
        String remotePath = tempRemote.remotes.get(remoteName);
        File remoteGitletFile = join(remotePath);
        if (!remoteGitletFile.exists()) {
            System.out.println("Remote directory not found.");
            System.exit(0);
        }
        File remoteHeadFile = join(remoteGitletFile, ".HEAD");
        String remoteHeadCommit = readContentsAsString(remoteHeadFile);
        CommitTree tempCommitTree = readObject(Repository.COMMIT_TREE, CommitTree.class);
        if (!tempCommitTree.commitTree.containsKey(remoteHeadCommit)) {
            System.out.println("Please pull down remote changes before pushing.");
            System.exit(0);
        }
        File remoteBranchFile = join(remoteGitletFile + File.separator + ".branches/",
                remoteBranchName);
        if (!remoteBranchFile.exists()) {

        }
    }

}