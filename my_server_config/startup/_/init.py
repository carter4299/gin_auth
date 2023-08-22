import subprocess
import urllib.request
import re
import time
import os 
import shutil

#python3 -m http.server 8000

def get_terminal_width():
    return shutil.get_terminal_size().columns
WIDTH = get_terminal_width()

ssh_setup = {
    "updated": "false",
    "server host": "?",
    "public_key": "?",
    "private_key": "?"
}

"""nginx = {
    "updated": "false",
    "domain_name": "?",
    "host": "?",
    "port": "?"
}"""

gmod = {
    "updated": "false",
    "github_username": "?",
    "github_repo": "?"
}

env = {
    "updated": "false",
    "server host": "?",
    "server port": "?",
    "web/proxy host": "?",
    "web/proxy port": "?",
}
"""-----------------------------------------------------------CLI-----------------------------------------------------------"""
def display_boxed_info():
    os.system('clear' if os.name == 'posix' else 'cls')
    dec_(" Entering Python Script ")
    dec("ssh setup")
    for key, value in ssh_setup.items():
        c_format(f"{key.capitalize()} = {value}")
    dec("go_mod.txt")
    for key, value in gmod.items():
        c_format(f"{key.capitalize()} = {value}")
    dec(".env")
    for key, value in env.items():
        c_format(f"{key.capitalize()} = {value}")
    dec_("*")

def c_format(s):
    print(f"* {s} "+ "*".rjust(WIDTH - len(s) - 3))

def dec(s):
    print("* " + (WIDTH - 4) * " " + " *" + "\n" + "* " + s.center(WIDTH - 4, ' ') + " *" + "\n" + "* " + (WIDTH - 4) * " " + " *")

def dec_(s):
    print("* " + s.center(WIDTH - 4, '*') + " *")

def dec__(s):
    print("* " + s.center(WIDTH - 4, ' ') + " *")

def update_env(key, value):
    if key in env:
        env[key] = value
    display_boxed_info()

def update_gmod(key, value):
    if key in gmod:
        gmod[key] = value
    display_boxed_info()

def update_ssh_setup(key, value):
    if key in ssh_setup:
        ssh_setup[key] = value
    display_boxed_info()

def write_env(_env):
    with open('./my_server_config/startup/_/.env', 'w') as f:
        for key, value in _env.items():
            f.write(f'{key}={value}\n')
    update_env('updated', "true")

def write_gmod(username, repo_name):
    with open('./my_server_config/startup/_/go_mod.txt', 'w') as f:
        f.write(f'github_username={username}\n')
        f.write(f'github_repo={repo_name}\n')
    update_gmod("updated", "true")

"""-----------------------------------------------------------USER INPUT-----------------------------------------------------------"""

def get_input_from_user():
    display_boxed_info() 

    _env = {}

    print("SSH Setup Script, only used to connect from your laptop etc... not really needed.")
    i = input("Would you like to setup an ssh server? (y/n): ")
    if i.lower() == 'y':
        _env['SSH_USER'], _env['SSH_SERVER'] = setup_ssh()
        update_ssh_setup("updated", "true")

    get_yaml()

    _env['SERVER_HOST'] = get_host('server', 'host')
    update_env('server host', _env['SERVER_HOST'])
    _env['SERVER_PORT'] = get_port('server', 'port')
    update_env('server port', _env['SERVER_PORT'])
    _env['WEB_PROXY_HOST'] = get_host('web/proxy', 'host')
    update_env('web/proxy host', _env['WEB_PROXY_HOST'])
    _env['WEB_PROXY_PORT'] = get_port('web/proxy', 'port')
    update_env('web/proxy port', _env['WEB_PROXY_PORT'])
    _env['AUTH_DB_USER'], _env['AUTH_DB_PASS'] = get_key(0, 0)
    if _env['AUTH_DB_USER'] == '' or _env['AUTH_DB_PASS'] == '':
        dec('ERROR: scraper is not working ...')
        exit()
    _env['AUTH_ENC_KEY'], _env['INDEX_KEY1'], _env['INDEX_KEY2'], _env['INDEX_KEY3'], _env['INDEX_KEY4'], _env['INDEX_KEY5'], _env['INDEX_KEY6'] = get_key(1, 0)
    if _env['AUTH_ENC_KEY'] == '' or _env['INDEX_KEY1'] == '' or _env['INDEX_KEY2'] == '' or _env['INDEX_KEY3'] == '' or _env['INDEX_KEY4'] == '' or _env['INDEX_KEY5'] == '' or _env['INDEX_KEY6'] == '':
        dec('ERROR: scraper is not working ...')
        exit()

    write_env(_env)


def get_yaml():
    print("Go mod values, do not have to be real, just valid.")
    username = input("Enter your github username: ")
    update_gmod("github_username", username)
    repo_name = input("Enter your github repo name: ")
    update_gmod("github_repo", repo_name)
    _i = input(f"Is this go mod correct -> github.com/{username}/{repo_name}? (y/n): ")
    if _i.lower() == "y":
        write_gmod(username, repo_name)
    elif _i.lower() == "n":
        print('... calling get_yaml() again ...')
        get_yaml()
    else:
        dec("ERROR: invalid input -> calling get_yaml() again ...")
        get_yaml()

def get_host(x, y):
    _p = input(f"Would you like to use localhost for your {x} {y}? (y/n): ")
    if _p.lower() == 'y':
        return 'localhost'
    elif _p.lower() == 'n':
        p = input(f'Please enter a host ip for your {x} {y}: ')
        p_ip = input(f"Is this host ip correct -> {p}? (y/n): ")
        if p_ip.lower() == 'y':
            return p
    else:
        c_format(f'... answer invalid {_p} -> calling get_host() again ...')
        get_host(x, y)

def get_port(x, z):
    if x == 'server' or x == 'server ssh tunnel':
        _p = input(f"Would you like to use Gin's Default port for your {x} {z}? (y/n): ")
        if _p.lower() == 'y':
            return '8080'
        if _p.lower() not in ['y', 'n']:
            print(f'... answer invalid {_p} -> calling get_port() again ...')
            get_port(x, z)
    p = input(f'Please enter a port for your {x} {z}:  ')
    p_ip = input(f"Is this port correct -> {p}? (y/n): ")
    if p_ip.lower() == 'y':
        return p
    else:
        print('... calling get_port() again ...')
        get_port(x, z)

def get_key(n, t):
    url = "https://generate-random.org/encryption-key-generator?count=7&bytes=16&cipher=aes-256-cbc&string=&password=" if n == 1 else "https://generate-random.org/encryption-key-generator?count=2&bytes=32&cipher=aes-256-cbc-hmac-sha256&string=&password="

    if t == 0:
        print("1st Attempt: Requesting keys from -> generate-random.org")
    elif t == 1:
        print("2nd Attempt: Waiting 10 seconds before requesting more key from -> generate-random.org")
        time.sleep(10)
    elif t == 2:
        print("3rd Attempt: Waiting 10 seconds before requesting more key from -> generate-random.org")
        time.sleep(10)
    else:
        print("\nERROR: scraper is not working ...\n")
        if n == 0:
            print("\t3 attempts have returned an error, in the set_env() function you can hard code the keys env['AUTH_DB_USER'], env['AUTH_DB_PASS'] ...")
            print("\t-> Please visit https://generate-random.org/encryption-key-generator?count=2&bytes=32&cipher=aes-256-cbc-hmac-sha256&string=&password= to get your keys ...")
            print("\t-> If the link is broken, they may have changed their website, please visit https://generate-random.org/ and find the key generator and use the following settings: Count=2, Bytes=32, Cipher=aes-256-cbc-hmac-sha256")
        if n == 1:
            print("\t3 attempts have returned an error, in the set_env() function you can hard code your keys env['SESSION_KEY'], env['KEY1'], env['KEY2'], env['KEY3'], env['KEY4'], env['KEY5'], env['KEY6'] ...")
            print("\t-> Please visit https://generate-random.org/encryption-key-generator?count=7&bytes=16&cipher=aes-256-cbc&string=&password= to get your keys ...")
            print("\t-> If the link is broken, they may have changed their website, please visit https://generate-random.org/ and find the key generator and use the following settings: Count=7, Bytes=16, Cipher=aes-256-cbc")
        exit()
    
    response = urllib.request.urlopen(url)
    if response.getcode() != 200:
        print(f"Error: {response.getcode()}")
        get_key(n, t + 1)

    html = response.read().decode('utf-8')

    ret = [k[1] for k in re.findall(r'<span class="text-warning monospace" title=(.*?)>(.*?)</span>', html)]
    print(f"\tKeys generated: {ret.__len__()}")

    if ret.__len__() == 7:
        return ret[0], ret[1], ret[2], ret[3], ret[4], ret[5], ret[6]
    if ret.__len__() == 2:
        return ret[0], ret[1]
    else:
        dec("ERROR: scraper is not returning correct amount of keys")

def setup_ssh():
    current_user = os.getlogin()
    server_address = get_server_address()
    
    if server_address:
        prompt_username = input(f"Enter your username (default: {current_user}): ")
        username = prompt_username or current_user
        promp_choice = input(f"Would you like to use (0)-{server_address} or (1)-localhost as the server address? (default: 0): ")
        if promp_choice == "1":
            server_address = 'localhost'
        print(f"\n{username}@{server_address}")
        correct = input("Is the above information correct? (y/n): ")
        
        if correct.lower() != 'y':
            dec("ERROR: invalid or 'n' input calling -> setup_ssh() again ...")
            setup_ssh()
    
    update_ssh_setup("server host", server_address)
    choice = input("Do you want to install the SSH server? (y/n): ")
    if choice.lower() == 'y':
        install_ssh_server()
    
    choice = input("Do you want to generate SSH keys for the current user? (y/n): ")
    if choice.lower() == 'y':
        print("\nSSH will now ask you what file you want to save the key pair in.\nPress enter to accept the default value.\nIt will then ask for a passphrase.\nYou can leave this blank if you want, though not recommended.")
        generate_ssh_keys()
        print(f"\nYour key pair is located at: ~/.ssh/id_rsa (private key) and ~/.ssh/id_rsa.pub (public key)")
    
    print("\nSetup complete!")

    return username, server_address

"""-----------------------------------------------------------INSTALLERS-----------------------------------------------------------"""

def install_ssh_server():
    """Install the SSH server on Ubuntu/Debian systems."""
    try:
        subprocess.run(["sudo", "apt-get", "update"], check=True)
        subprocess.run(["sudo", "apt-get", "install", "-y", "openssh-server"], check=True)
        print("SSH server installed successfully.")
    except subprocess.CalledProcessError:
        print("There was an error installing the SSH server. Are you running this on Ubuntu or another Debian-based system?")

def generate_ssh_keys():
    """Generate SSH keys for the current user."""
    if not os.path.exists(os.path.expanduser("~/.ssh/id_rsa")):
        try:
            subprocess.run(["ssh-keygen"], check=True)
            print("SSH keys generated successfully.")
        except subprocess.CalledProcessError:
            print("There was an error generating SSH keys.")
    else:
        print("SSH keys already exist for this user.")
    
    update_ssh_setup("public_key", os.path.expanduser("~/.ssh/id_rsa.pub"))
    update_ssh_setup("private_key", os.path.expanduser("~/.ssh/id_rsa"))

def get_server_address():
    """Retrieve the local IP address of the server."""
    try:
        ip_address = subprocess.getoutput("hostname -I").split()[0]
        return ip_address
    except:
        print("There was an error retrieving the server's IP address.")
        return None


if __name__ == '__main__':
    get_input_from_user()


"""
find . -type f -exec sed -i "s||$github_username|g" {} +
find . -type f -exec sed -i "s||$github_repo|g" {} +
"""