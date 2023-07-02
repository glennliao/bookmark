import os, shutil


if __name__ == '__main__':
    for item in os.walk("./temp"):
        if len(item[2]) > 0:
            for f in item[2]:
                to = os.path.join("./temp", f+"_"+os.path.basename(item[0]))
                if ".exe" in f:
                    to = to.replace(".exe","")+".exe"
                shutil.copy(
                    os.path.join(item[0], f),
                    to)
                shutil.rmtree(item[0])