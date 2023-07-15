import os, shutil
import tarfile,zipfile

def main():
    for item in os.walk("./temp"):
        if len(item[2]) > 0:
           for f in item[2]:
                if ".exe" in f:
                    tar_file = "./temp/"+os.path.basename(item[0])+".zip"
                    zip = zipfile.ZipFile(tar_file, mode='w')
                    zip.write(os.path.join(item[0], f), arcname=f,compress_type=zipfile.ZIP_DEFLATED)
                    zip.close()

                else:
                    tar_file = "./temp/"+os.path.basename(item[0])+".tar.gz"
                    tar = tarfile.open(tar_file, mode='w:gz')
                    tar.add(os.path.join(item[0], f), arcname=f)
                    tar.close()

           shutil.rmtree(item[0])

if __name__ == '__main__':
    main()