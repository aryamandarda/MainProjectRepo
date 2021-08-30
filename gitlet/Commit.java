package gitlet;

import java.io.Serializable;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.HashMap;

/** Represents a gitlet commit object.
 *  Provides the instance variables, methods, etc. for Commit Objects
 *  does at a high level.
 *
 *  @author Aryaman Darda
 */
public class Commit implements Serializable {
    /**
     * List all instance variables of the Commit class here with a useful
     * comment above them describing what that variable represents and how that
     * variable is used. We've provided one example for `message`.
     */

    /** The message of this Commit. */
    protected String message;
    /** The date and time of the Commit */
    protected LocalDateTime date;
    protected String hashName;
    protected ArrayList<String> parentCommit;
    protected HashMap<String, String> blobMap;
    protected String branch;

    public Commit() {
        this.message = "initial commit";
        this.date = LocalDateTime.parse("1970-01-01T00:00:00");
        this.hashName = Utils.sha1(Utils.serialize(this));
        this.parentCommit = null;
        this.blobMap = new HashMap<>();
        this.branch = "master";
    }
    public Commit(String message, LocalDateTime dateTime) {
        this.message = message;
        this.date = dateTime;
        this.parentCommit = new ArrayList<>();
        this.parentCommit.add(Utils.readContentsAsString(Repository.HEAD));
        this.blobMap = new HashMap<>();
        this.branch = Utils.readContentsAsString(Repository.CURRENT_BRANCH);
    }

    public String getHashName() {
        return this.hashName;
    }

    public void setHashName() {
        this.hashName = Utils.sha1(Utils.serialize(this));
    }

    public void printFormattedDate() {
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("EEE MMM dd HH:mm:ss yyyy Z");
        ZoneId zoneId = ZoneId.of("US/Pacific");
        ZonedDateTime zonedDateTime = this.date.atZone(zoneId);
        System.out.println("Date: " + zonedDateTime.format(formatter));
    }

}
