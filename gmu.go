package main

import (
	"flag"
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"log"
	"gopkg.in/ini.v1"
	"github.com/ysouyno/gmu/utils"
)

const VERSION = "0.0.1"
const GMU_CONFIG = ".gmuconfig"
const GITCONFIG = ".gitconfig"
const SSHCONFIG = ".ssh"
const SEC_NAME_CUR = "current"
const SEC_NAME_USERS = "users"
const KEY_NAME_GIT = "gitconfig"
const KEY_NAME_SSH = "sshconfig"
const KEY_NAME_USER = "user"
const KEY_NAME_NAME = "name"
const KEY_NAME_EMAIL = "email"

var flag_ver *bool
var flag_cur *bool
var flag_upd *bool
var flag_all *bool
var flag_chk string

func get_git_config_info() bool {
	USERPROFILE := os.Getenv(`USERPROFILE`)

	gitconfig := USERPROFILE + "\\" + GITCONFIG

	if !utils.FileExist(gitconfig) {
		log.Printf("No %s found.\n", GITCONFIG)
		return false
	}

	cfg, err := ini.Load(gitconfig)
	if err != nil {
		log.Println(err)
		return false
	}

	name := cfg.Section(KEY_NAME_USER).Key(KEY_NAME_NAME).String()
	email := cfg.Section(KEY_NAME_USER).Key(KEY_NAME_EMAIL).String()

	fmt.Printf("Current user: %s <%s>\n", name, email)

	return true
}

func get_current_git_user() string {
	USERPROFILE := os.Getenv(`USERPROFILE`)

	gitconfig := USERPROFILE + "\\" + GITCONFIG

	if !utils.FileExist(gitconfig) {
		log.Printf("No %s found.\n", GITCONFIG)
		return ""
	}

	cfg, err := ini.Load(gitconfig)
	if err != nil {
		log.Println(err)
		return ""
	}

	return cfg.Section(KEY_NAME_USER).Key(KEY_NAME_NAME).String()
}

func update_ini_config(home string) bool {
	ini_config := home + "\\" + GMU_CONFIG;
	if !utils.FileExist(ini_config) {
		file, err := os.Create(ini_config)
		if err != nil {
			log.Println(err)
			return false
		}
		defer file.Close()
	}

	cfg, err := ini.Load(ini_config)
	if err != nil {
		log.Println(err)
		return false
	}

	// handle [current]
	curr_sec := cfg.Section(SEC_NAME_CUR)
	if !curr_sec.HasKey(KEY_NAME_NAME) {
		curr_sec.NewKey(KEY_NAME_NAME, "");
	}

	if !curr_sec.HasKey(KEY_NAME_GIT) {
		curr_sec.NewKey(KEY_NAME_GIT, home + "\\" + GITCONFIG)
	}

	if !curr_sec.HasKey(KEY_NAME_SSH) {
		curr_sec.NewKey(KEY_NAME_SSH, home + "\\" + SSHCONFIG)
	}

	// update current git user
	user := get_current_git_user()
	curr_sec.Key(KEY_NAME_NAME).SetValue(user)

	// handle [users]
	users_sec := cfg.Section(SEC_NAME_USERS)
	if !users_sec.HasKey(KEY_NAME_NAME) {
		users_sec.NewKey(KEY_NAME_NAME, user)
	}

	users := users_sec.Key(KEY_NAME_NAME).String()
	if !strings.Contains(users, user) {
		// new git user
		users = users + " " + user
		users_sec.Key(KEY_NAME_NAME).SetValue(users)
	}

	// handle [%git user%]
	new_user_sec := cfg.Section(user)

	gitconfig := home + "\\" + GITCONFIG + "." + user
	if utils.FileExist(gitconfig) {
		new_user_sec.NewKey(KEY_NAME_GIT, gitconfig)
	}

	sshconfig := home + "\\" + SSHCONFIG + "." + user
	if utils.FileExist(sshconfig) {
		new_user_sec.NewKey(KEY_NAME_SSH, sshconfig)
	}

	cfg.SaveTo(ini_config)

	return true
}

func save_git_config(home, user string) bool {
	current_git_user := get_current_git_user()
	old_config_file := home + "\\" + GITCONFIG
	new_config_file := old_config_file + "." + user

	// if current user is equal to 'user' and .gitconfig.user exists,
	// no need save again
	if current_git_user == user && utils.FileExist(new_config_file) {
		return true
	}

	// save .gitconfig as .gitconfig.user_name
	ret, _ := utils.CopyFile(new_config_file, old_config_file)
	if ret == 0 {
		return false
	}

	return true
}

func save_ssh_config(home, user string) bool {
	current_git_user := get_current_git_user()
	old_config_file := home + "\\" + SSHCONFIG
	new_config_file := old_config_file + "." + user

	// if current user is equal to 'user' and .ssh.user exists,
	// no need save again
	if current_git_user == user && utils.FileExist(new_config_file) {
		return true
	}

	err := os.MkdirAll(new_config_file, os.ModePerm)
	if err != nil {
		log.Println(err)
		return false
	}

	files, err := ioutil.ReadDir(old_config_file)
	if err != nil {
		log.Println(err)
		return false
	}

	for _, f := range files {
		f_old := old_config_file + "\\" + f.Name()
		f_new := new_config_file + "\\" + f.Name()
		utils.CopyFile(f_new, f_old)
	}

	return true
}

func init_env() bool {
	USERPROFILE := os.Getenv(`USERPROFILE`)

	gitconfig := USERPROFILE + "\\" + GITCONFIG

	if !utils.FileExist(gitconfig) {
		log.Printf("No %s found.\n", GITCONFIG)
		return false
	}

	sshconfig := USERPROFILE + "\\" + SSHCONFIG
	if !utils.FileExist(sshconfig) {
		log.Printf("No %s found.\n", SSHCONFIG)
		return false
	}

	cfg, err := ini.Load(gitconfig)
	if err != nil {
		log.Println(err)
		return false
	}

	name := cfg.Section(KEY_NAME_USER).Key(KEY_NAME_NAME).String()

	if !save_git_config(USERPROFILE, name) {
		return false
	}

	if !save_ssh_config(USERPROFILE, name) {
		return false
	}

	update_ini_config(USERPROFILE)

	return true
}

func update_env() bool {
	return init_env()
}

func list_user() bool {
	home := os.Getenv(`USERPROFILE`)

	ini_config := home + "\\" + GMU_CONFIG;
	if !utils.FileExist(ini_config) {
		log.Printf("No %s found.\n", GMU_CONFIG)
		return false
	}

	cfg, err := ini.Load(ini_config)
	if err != nil {
		log.Println(err)
		return false
	}

	curr_user := cfg.Section(SEC_NAME_CUR).Key(KEY_NAME_NAME).String()
	users := cfg.Section(SEC_NAME_USERS).Key(KEY_NAME_NAME).String()

	user_arr := strings.Fields(users)
	for _, ele := range user_arr {
		if ele == curr_user {
			fmt.Println("*", ele)
		} else {
			fmt.Println(" ", ele)
		}
	}

	return true
}

func checkout_user(user string) bool {
	home := os.Getenv(`USERPROFILE`)

	ini_config := home + "\\" + GMU_CONFIG;
	if !utils.FileExist(ini_config) {
		log.Printf("No %s found.\n", GMU_CONFIG)
		return false
	}

	cfg, err := ini.Load(ini_config)
	if err != nil {
		log.Println(err)
		return false
	}

	curr_user := cfg.Section(SEC_NAME_CUR).Key(KEY_NAME_NAME).String()
	if curr_user == user {
		fmt.Printf("%s is already the current user.\n", user)
		return true
	}

	users := cfg.Section(SEC_NAME_USERS).Key(KEY_NAME_NAME).String()
	if !strings.Contains(users, user) {
		fmt.Printf("%s does not exist, can't checkout.", user)
		return false
	}

	user_gitconfig := cfg.Section(user).Key(KEY_NAME_GIT).String()
	if !utils.FileExist(user_gitconfig) {
		log.Printf("not find %s's %s", user, GITCONFIG)
		return false
	}

	curr_gitconfig := cfg.Section(SEC_NAME_CUR).Key(KEY_NAME_GIT).String()
	ret, _ := utils.CopyFile(curr_gitconfig, user_gitconfig)
	if ret == 0 {
		log.Printf("checkout %s failed\n", GITCONFIG)
		return false
	}

	user_sshconfig := cfg.Section(user).Key(KEY_NAME_SSH).String()
	if !utils.FileExist(user_sshconfig) {
		log.Printf("not find %s's %s", user, SSHCONFIG)
		return false
	}

	curr_sshconfig := cfg.Section(SEC_NAME_CUR).Key(KEY_NAME_SSH).String()

	files, err := ioutil.ReadDir(curr_sshconfig)
	if err != nil {
		log.Println(err)
		return false
	}

	for _, f := range files {
		f_old := user_sshconfig + "\\" + f.Name()
		f_new := curr_sshconfig + "\\" + f.Name()

		ret, _ = utils.CopyFile(f_new, f_old)
		if ret == 0 {
			log.Printf("checkout %s failed", SSHCONFIG)
			return false
		}
	}

	update_env()
	list_user()

	return true
}

func init() {
	log.SetPrefix("[gmu] ")
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	flag_ver = flag.Bool("v", false, "Print the version number.")
	flag_cur = flag.Bool("i", false, "Show current git user <email>.")
	flag_upd = flag.Bool("u", false, "Update gmu.")
	flag_all = flag.Bool("a", false, "Print all git users.")
	flag.StringVar(&flag_chk, "c", "", "Set `user` as the current user.")

	init_env()
}

func main() {
	flag.Parse()

	if *flag_ver {
		fmt.Println("Version:", VERSION)
	} else if *flag_cur {
		get_git_config_info()
	} else if *flag_upd {
		update_env()
	} else if *flag_all {
		list_user()
	} else if flag_chk != "" {
		checkout_user(flag_chk)
	} else {
		fmt.Println("Try 'gmu -h' for more options.")
	}
}
