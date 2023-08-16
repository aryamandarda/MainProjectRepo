# You may want to install "gprof2dot"
import io
from collections import Counter
import csv
import numpy as np
import scipy.io
import sklearn.utils
import sklearn.model_selection
import sklearn.tree
import sklearn.metrics
from numpy import genfromtxt
from scipy import stats
from sklearn.base import BaseEstimator, ClassifierMixin
import graphviz
import pydot
import gprof2dot
import matplotlib.pyplot as plt

np.random.seed(6996) # Random seed for reproducibility

eps = 1e-5  # a small number


class DecisionTree:
    def __init__(self, max_depth=3, feature_labels=None, m=None):
        self.max_depth = max_depth
        self.features = feature_labels
        self.left, self.right = None, None  # for non-leaf nodes
        self.split_idx, self.thresh = None, None  # for non-leaf nodes
        self.data, self.pred = None, None  # for leaf nodes
        self.m = m

    @staticmethod
    def information_gain(X, y, thresh):
        lsplit, rsplit = np.where(X < thresh)[0], np.where(X >= thresh)[0]
        yl, yr = y[lsplit], y[rsplit]

        def entropy(y):
            counts = np.unique(y, return_counts=True)[1]
            n = len(y)
            sum = 0
            for c in counts:
                prob = c/n
                sum -= prob*np.log2(prob)
            return sum
        
        Hl, Hr = entropy(yl), entropy(yr)
        return -1 * ((len(yl) * Hl + len(yr) * Hr) / (len(yl) + len(yr)))

    @staticmethod
    def gini_impurity(X, y, thresh):
        # TODO: implement gini impurity function
        pass

    def split(self, X, y, idx, thresh):
        X0, idx0, X1, idx1 = self.split_test(X, idx=idx, thresh=thresh)
        y0, y1 = y[idx0], y[idx1]
        return X0, y0, X1, y1

    def split_test(self, X, idx, thresh):
        idx0 = np.where(X[:, idx] < thresh)[0]
        idx1 = np.where(X[:, idx] >= thresh)[0]
        X0, X1 = X[idx0, :], X[idx1, :]
        return X0, idx0, X1, idx1

    def fit(self, X, y):
        if self.max_depth > 0:
            # compute entropy gain for all single-dimension splits,
            # thresholding with a linear interpolation of 10 values
            gains = []
            # The following logic prevents thresholding on exactly the minimum
            # or maximum values, which may not lead to any meaningful node
            # splits.

            # Randomly choosing attributes for fitting data to DT
            randomAttributes = None
            if self.m:
                randomAttributes = np.random.choice(np.arange(len(self.features)), size=self.m, replace=False)

            thresh = np.array([
                np.linspace(np.min(X[:, i]) + eps, np.max(X[:, i]) - eps, num=10)
                for i in range(X.shape[1])
            ])
            if self.m:
                for i in range(X.shape[1]):
                    if i in randomAttributes:
                        gains.append([self.information_gain(X[:, i], y, t) for t in thresh[i, :]])
                    else:
                        gains.append([-1000 for t in thresh[i, :]])
            else:
                for i in range(X.shape[1]):
                    gains.append([self.information_gain(X[:, i], y, t) for t in thresh[i, :]])
            gains = np.nan_to_num(np.array(gains))
            self.split_idx, thresh_idx = np.unravel_index(np.argmax(gains), gains.shape)
            self.thresh = thresh[self.split_idx, thresh_idx]
            X0, y0, X1, y1 = self.split(X, y, idx=self.split_idx, thresh=self.thresh)
            if X0.size > 0 and X1.size > 0:
                self.left = DecisionTree(
                    max_depth=self.max_depth - 1, feature_labels=self.features)
                self.left.fit(X0, y0)
                self.right = DecisionTree(
                    max_depth=self.max_depth - 1, feature_labels=self.features)
                self.right.fit(X1, y1)
            else:
                self.max_depth = 0
                self.data, self.labels = X, y
                self.pred = stats.mode(y, keepdims=True).mode[0]
        else:
            self.data, self.labels = X, y
            self.pred = stats.mode(y, keepdims=True).mode[0]
        return self

    def predict(self, X):
        if self.max_depth == 0:
            return self.pred * np.ones(X.shape[0])
        else:
            X0, idx0, X1, idx1 = self.split_test(X, idx=self.split_idx, thresh=self.thresh)
            yhat = np.zeros(X.shape[0])
            yhat[idx0] = self.left.predict(X0)
            yhat[idx1] = self.right.predict(X1)
            return yhat

    def __repr__(self):
        if self.max_depth == 0:
            return "%s (%s)" % (self.pred, self.labels.size)
        else:
            return "[%s < %s: %s | %s]" % (self.features[self.split_idx],
                                           self.thresh, self.left.__repr__(),
                                           self.right.__repr__())


class BaggedTrees(BaseEstimator, ClassifierMixin):
    def __init__(self, params=None, n=200, sampleSize=None, maxDepth=3, features=None):
        if params is None:
            params = {}
        self.params = params
        self.n = n
        self.sampleSize = sampleSize
        self.maxDepth = maxDepth
        self.features = features 
        self.decision_trees = [DecisionTree(max_depth=self.maxDepth, feature_labels=features, m=None) for i in range(self.n)]

    def fit(self, X, y):
        assert self.sampleSize <= len(y)
        for tree in self.decision_trees:
            randomSamples = np.random.choice(np.arange(len(y)), size=self.sampleSize, replace=True)
            trainX = X[randomSamples, :]
            trainy = y[randomSamples]

            tree.fit(trainX, trainy)

    def predict(self, X):
        preds = []
        for tree in self.decision_trees:
            preds.append(tree.predict(X))
        predArr = np.vstack(preds)
        predModeArr = stats.mode(predArr, keepdims=True)
        predMode = stats.mode(predModeArr, keepdims=True)[0].reshape(-1)
        return predMode

class RandomForest(BaggedTrees):
    def __init__(self, params=None, n=200, m=1, sampleSize=None, maxDepth=1, features=None):
        if params is None:
            params = {}
        self.params = params
        self.n = n
        self.m = m
        self.sampleSize = sampleSize
        self.maxDepth = maxDepth
        self.features = features
        self.decision_trees = [DecisionTree(max_depth=self.maxDepth, feature_labels=self.features, m=self.m) for i in range(self.n)]


class BoostedRandomForest(RandomForest):
    def fit(self, X, y):
        self.w = np.ones(X.shape[0]) / X.shape[0]  # Weights on data
        self.a = np.zeros(self.n)  # Weights on decision trees
        # TODO: implement function
        return self

    def predict(self, X):
        # TODO: implement function
        pass


def preprocess(data, fill_mode=True, min_freq=10, onehot_cols=[]):
    # fill_mode = False

    # Temporarily assign -1 to missing data
    data[data == ''] = '-1'

    # Hash the columns (used for handling strings)
    onehot_encoding = []
    onehot_features = []
    for col in onehot_cols:
        counter = Counter(data[:, col])
        for term in counter.most_common():
            if term[0] == '-1':
                continue
            if term[-1] <= min_freq:
                break
            onehot_features.append(term[0])
            onehot_encoding.append((data[:, col] == term[0]).astype(float))
        data[:, col] = '0'
    onehot_encoding = np.array(onehot_encoding).T
    data = np.hstack([np.array(data, dtype=float), np.array(onehot_encoding)])

    # Replace missing data with the mode value. We use the mode instead of
    # the mean or median because this makes more sense for categorical
    # features such as gender or cabin type, which are not ordered.
    if fill_mode:
        for i in range(data.shape[-1]):
            mode = stats.mode(data[((data[:, i] < -1 - eps) +
                                    (data[:, i] > -1 + eps))][:, i], keepdims=True).mode[0]
            data[(data[:, i] > -1 - eps) * (data[:, i] < -1 + eps)][:, i] = mode

    return data, onehot_features


def evaluate(clf):
    print("Cross validation", sklearn.model_selection.cross_val_score(clf, X, y))
    if hasattr(clf, "decision_trees"):
        counter = Counter([t.tree_.feature[0] for t in clf.decision_trees])
        first_splits = [(features[term[0]], term[1]) for term in counter.most_common()]
        print("First splits", first_splits)


if __name__ == "__main__":
    # Results for the SPAM dataset
    dataset = "spam"
    features = [
            "pain", "private", "bank", "money", "drug", "spam", "prescription", "creative",
            "height", "featured", "differ", "width", "other", "energy", "business", "message",
            "volumes", "revision", "path", "meter", "memo", "planning", "pleased", "record", "out",
            "semicolon", "dollar", "sharp", "exclamation", "parenthesis", "square_bracket",
            "ampersand"
        ]
    assert len(features) == 32

    # Load spam data
    path_train = './hw5_code/dataset/spam/spam_data.mat'
    data = scipy.io.loadmat(path_train)
    X = data['training_data']
    y = np.squeeze(data['training_labels'])

    shuffled = sklearn.utils.shuffle(X, y, random_state=6996)
    trainSize = int(np.round(X.shape[0] * 0.8))
    XShuff, yShuff = shuffled[0], shuffled[1]
    spamTrainData, spamTrainLabels = XShuff[:trainSize], yShuff[:trainSize]
    spamValData, spamValLabels = XShuff[trainSize:], yShuff[trainSize:]

    Z = data['test_data']
    class_names = ["Ham", "Spam"]

    # Fitting to Decision Tree and getting accuracy of train & validation data. Predicting on test.
    tree = DecisionTree(max_depth=28, feature_labels=features, m=None)
    tree.fit(spamTrainData, spamTrainLabels)
    
    trainPred = tree.predict(spamTrainData)
    trainScore = sklearn.metrics.accuracy_score(spamTrainLabels, trainPred)
    print("Decision Tree Accuracy Score for SPAM train data: ", trainScore)

    valPred = tree.predict(spamValData)
    valScore = sklearn.metrics.accuracy_score(spamValLabels, valPred)
    print("Decision Tree Accuracy Score for SPAM validation data: ", valScore)

    # Fitting to Random Forest and getting accuracy of train & validation data
    forest = RandomForest(params=None, n=100, m=6, sampleSize=4000, maxDepth=11, features=features)
    forest.fit(spamTrainData, spamTrainLabels)

    # Kaggle
    spamResults = forest.predict(Z)
    with open('spam_hw5', mode='w') as file:
        print("starting...")
        writer = csv.writer(file, delimiter=',')
        writer.writerow(['Id', 'Category'])
        for idx, x in enumerate(spamResults):
            writer.writerow([idx+1, int(x)])

    trainPred = forest.predict(spamTrainData)
    trainScore = sklearn.metrics.accuracy_score(spamTrainLabels, trainPred)
    print("Random Forest Accuracy Score for SPAM train data: ", trainScore)

    valPred = forest.predict(spamValData)
    valScore = sklearn.metrics.accuracy_score(spamValLabels, valPred)
    print("Random Forest Accuracy Score for SPAM validation data: ", valScore)

    # Split Visualization
    def visualizeSplit(treeboi, data, classNames):
        if not treeboi.pred is None:
            print("Therefore, this email was ",classNames[treeboi.pred])
            return
        else:
            output = "('" + str(treeboi.features[treeboi.split_idx])
            if data[treeboi.split_idx] < treeboi.thresh:
                output += "') <" + str(treeboi.thresh)
                print(output)
                visualizeSplit(treeboi.left, data, classNames)
            else:
                output += "') >=" + str(treeboi.thresh)
                print(output)
                visualizeSplit(treeboi.right, data, classNames)

    visualizeSplit(tree, spamTrainData[0], class_names)
    print(" ")
    visualizeSplit(tree, spamTrainData[14], class_names)
    
    # Tree depth hyperparameter tuning
    depthRange = np.arange(1, 46)
    accuracy = []
    for depth in depthRange:
        dtree = DecisionTree(max_depth=depth, feature_labels=features, m=None)
        dtree.fit(spamTrainData, spamTrainLabels)
        valP = dtree.predict(spamValData)
        accuracy.append(sklearn.metrics.accuracy_score(spamValLabels, valP))

    # Load Titanic data
    dataset = "titanic"
    path_train = './hw5_code/dataset/titanic/titanic_training.csv'
    data = genfromtxt(path_train, delimiter=',', dtype=None, encoding=None)
    path_test = './hw5_code/dataset/titanic/titanic_test_data.csv'
    test_data = genfromtxt(path_test, delimiter=',', dtype=None, encoding=None)
    y = data[1:, -1]  # label = survived
    class_names = ["Died", "Survived"]
    labeled_idx = np.where(y != '')[0]

    y = np.array(y[labeled_idx])
    y = y.astype(float).astype(int)

    print("\n\nPart (b): preprocessing the titanic dataset")
    X, onehot_features = preprocess(data[1:, :-1], onehot_cols=[1, 5, 7, 8])
    X = X[labeled_idx, :]
    Z, _ = preprocess(test_data[1:, :], onehot_cols=[1, 5, 7, 8])
    assert X.shape[1] == Z.shape[1]
    features = list(data[0, :-1]) + onehot_features

    shuffled = sklearn.utils.shuffle(X, y, random_state=6996)
    trainSize = int(np.round(X.shape[0] * 0.8))
    XShuff, yShuff = shuffled[0], shuffled[1]
    titanicTrainData, titanicTrainLabels = XShuff[:trainSize], yShuff[:trainSize]
    titanicValData, titanicValLabels = XShuff[trainSize:], yShuff[trainSize:]

    # Fitting to Decision Tree and getting accuracy of train & validation data
    tree = DecisionTree(max_depth=11, feature_labels=features, m=None)
    tree.fit(titanicTrainData, titanicTrainLabels)

    trainPred = tree.predict(titanicTrainData)
    trainScore = sklearn.metrics.accuracy_score(titanicTrainLabels, trainPred)
    print("Decision Tree Accuracy Score for Titanic train data: ", trainScore)

    valPred = tree.predict(titanicValData)
    valScore = sklearn.metrics.accuracy_score(titanicValLabels, valPred)
    print("Decision Tree Accuracy Score for Titanic validation data: ", valScore)

    # Fitting to Random Forest and getting accuracy of train & validation data
    forest = RandomForest(params=None, n=100, m=3, sampleSize=600, maxDepth=8, features=features)
    forest.fit(titanicTrainData, titanicTrainLabels)

    trainPred = forest.predict(titanicTrainData)
    trainScore = sklearn.metrics.accuracy_score(titanicTrainLabels, trainPred)
    print("Random Forest Accuracy Score for Titanic train data: ", trainScore)

    valPred = forest.predict(titanicValData)
    valScore = sklearn.metrics.accuracy_score(titanicValLabels, valPred)
    print("Random Forest Accuracy Score for Titanic validation data: ", valScore)

    # Visualizing depth 3 Tree
    dtree = DecisionTree(max_depth=3, feature_labels=features, m=None)
    dtree.fit(titanicTrainData, titanicTrainLabels)

    def visualizeTree(node, classNames, depth=0):
        output = " " * 10 * depth
        if not node.pred is None:
            print(output + classNames[node.pred])
            return
        else:
            output += node.features[node.split_idx]
            print(output + " >= " + str(node.thresh))
            visualizeTree(node.right, classNames, depth + 1)
            print(output + " < " + str(node.thresh))
            visualizeTree(node.left, classNames, depth + 1)

    visualizeTree(dtree, class_names)

    # Kaggle
    titanicResults = forest.predict(Z)
    with open('titanic_hw5', mode='w') as file:
        writer = csv.writer(file, delimiter=',')
        writer.writerow(['Id', 'Category'])
        for idx, x in enumerate(titanicResults):
            writer.writerow([idx+1, int(x)])
    
    # Graphing Hyperparameter Tuning results
    plt.plot(depthRange, accuracy)
    plt.xlabel("Max depth values")
    plt.ylabel("SPAM validation accuracy")
    plt.title("Max depth vs Validation accuracy for SPAM")
    plt.show()
