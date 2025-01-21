import os.path
import re
import shutil
import subprocess

from app.cmd import delete_file_if_exists, delete_dir_if_exists


class Builder(object):
    def __init__(self, build_dir: str):
        self.build_dir = build_dir
        self.lib_dir = os.path.join(self.build_dir, "..")

    def clean_lib_files(self, files: list[str]):
        for file in files:
            file_path = os.path.join(self.lib_dir, file)
            delete_file_if_exists(file_path)

    def clean_lib_dirs(self, dirs: list[str]):
        for dir_name in dirs:
            dir_path = os.path.join(self.lib_dir, dir_name)
            delete_dir_if_exists(dir_path)

    def append_lib_file_with_lines(self, file_name: str, lines: list[str]):
        file_path = os.path.join(self.lib_dir, file_name)
        with open(file_path, "r") as f:
            new_lines = f.readlines()
        new_lines.extend(lines)
        with open(file_path, "w") as f:
            f.writelines(new_lines)

    def prepare_go(self):
        clean_files = ["go.mod", "go.sum"]
        self.clean_lib_files(clean_files)
        os.chdir(self.lib_dir)
        ret = subprocess.run(["go", "mod", "init", "github.com/sagernet/sing-box/experimental/libbox/libXray"])
        if ret.returncode != 0:
            raise Exception("go mod init failed")
        self.append_lib_file_with_lines(
            "go.mod",
            [
                "\nreplace github.com/xtls/xray-core => ../Xray-core\n",
            ],
        )
        ret = subprocess.run(["go", "mod", "tidy"])
        if ret.returncode != 0:
            raise Exception("go mod tidy failed")

    def fix_go_module_version(self):
        # go module major version can NOT support v24, wtf?
        file_path = "./go.mod"
        with open(file_path, mode="r") as f:
            lines = f.readlines()
            lines.append("replace github.com/xtls/xray-core => ../Xray-core\n")

        with open(file_path, mode="w") as f:
            f.writelines(lines)

    def download_geo(self):
        os.chdir(self.lib_dir)
        main_path = os.path.join("main", "main.go")
        ret = subprocess.run(["go", "run", main_path])
        if ret.returncode != 0:
            raise Exception("download_geo failed")

    def prepare_gomobile(self):
        ret = subprocess.run(
            ["go", "install", "golang.org/x/mobile/cmd/gomobile@latest"]
        )
        if ret.returncode != 0:
            raise Exception("go install gomobile failed")
        ret = subprocess.run(["gomobile", "init"])
        if ret.returncode != 0:
            raise Exception("gomobile init failed")
        ret = subprocess.run(["go", "get", "golang.org/x/mobile/cmd/gomobile"])
        if ret.returncode != 0:
            raise Exception("gomobile update failed")
        ret = subprocess.run(["go", "get", "google.golang.org/genproto"])
        if ret.returncode != 0:
            raise Exception("gomobile install genproto failed")

    def prepare_static_lib(self):
        self.copy_template_file()
        self.fix_package_name()

    def copy_template_file(self):
        files = ['main.go']
        for file in files:
            src_file = os.path.join(self.build_dir, "template", file)
            shutil.copy(src_file, self.lib_dir)

    def fix_package_name(self):
        files = os.listdir(self.lib_dir)
        for file in files:
            if file.endswith(".go"):
                self.replace_package_name(file)

    def replace_package_name(self, file_name: str):
        file_path = os.path.join(self.lib_dir, file_name)
        new_lines = []
        with open(file_path, "r") as f:
            lines = f.readlines()
            for line in lines:
                new_line = line
                if re.match(r"^package\s+libXray", line):
                    new_line = "package main\n"
                new_lines.append(new_line)
        with open(file_path, "w") as f:
            f.writelines(new_lines)

    def before_build(self):
        self.prepare_go()
        self.download_geo()

    def build(self):
        pass

    def after_build(self):
        pass

    def reset_files(self):
        self.clean_lib_files(["main.go"])
        files = os.listdir(self.lib_dir)
        for file in files:
            if file.endswith(".go"):
                self.reset_package_name(file)

    def reset_package_name(self, file_name: str):
        file_path = os.path.join(self.lib_dir, file_name)
        new_lines = []
        with open(file_path, "r") as f:
            lines = f.readlines()
            for line in lines:
                new_line = line
                if re.match(r"^package\s+main", line):
                    new_line = "package libXray\n"
                new_lines.append(new_line)
        with open(file_path, "w") as f:
            f.writelines(new_lines)