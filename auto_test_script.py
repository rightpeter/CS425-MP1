import re
import os
import string
import random
import json


# reference: http://www.bswen.com/2018/04/python-How-to-generate-random-large-file-using-python.html
def generate_random_file(filename, size, know_pattern):
    random_chars = [random.choice(string.letters) for i in range(size)]
    random_ind = random.randint(0, size)
    # Insert new lines randomly
    i = 0
    while i < size:
        i += random.randint(10, 30)
        random_chars.insert(i, '\n')
    # Insert known pattern
    random_chars.insert(random_ind, '\n' + know_pattern + '\n')
    with open(filename, 'w+') as f:
        f.write(''.join(random_chars))


def send_log_file(vm_name, file_name):
    os.system("scp ./{0} {1}:/tmp/".format(file_name, vm_name))


def edit_config_file(filepath, vm_num):
    config = None
    with open(filepath, 'r') as f:
        config = json.load(f)
    config['current']['id'] = vm_num
    config['current']['log_path'] = '/tmp/vm{0}.test.log'.format(str(vm_num))
    with open(filepath, 'w') as f:
        json.dump(config, f)


def send_config_file(vm):
    edit_config_file('./mp1.test.config.json', int(vm.split('vm')[-1]))
    os.system("scp ./mp1.test.config.json {0}:/tmp/".format(vm))


# Will send config file and log file to all the vms
def send_files_to_vm(vm):
    log_file_name = "{0}.test.log".format(vm)
    known_pattern = "qwerty{0}qwerty\nzxcv{0}zxcv".format(vm)
    generate_random_file(log_file_name, 1024, known_pattern)
    send_log_file(vm, log_file_name)
    send_config_file(vm)


def get_line_count(output):
    return re.findall("VM([0-9]+): ([a-zA-Z0-9]+)", output)


def unit_test(grep_args, expected_line_count):
    output = os.popen('./client.bin "{0}"'.format(grep_args)).read()
    line_count = get_line_count(output)
    print("Testing pattern: {0}".format(grep_args))
    for v in line_count:
        if v[1] != expected_line_count:
            print(
                "Unit test failed for arguments: {} for VM{}. Expected line count: {} Got: {}".
                format(grep_args, v[0], expected_line_count, v[1]))
        else:
            print(
                "Unit test passed for arguments: {} for VM{}!!".
                format(grep_args, v[0]))


def shutdown_normal_servers_and_start_test_servers():
    """Shutdown normal servers"""
    num_vms = 10
    for i in range(1, num_vms + 1):
        print('Shutdown normal server for vm: ', i)
        os.system("ssh vm{} supervisorctl stop MP1".format(i))
        os.system("ssh vm{} supervisorctl start MP1_Test".format(i))


def shutdown_test_servers_and_start_normal_servers():
    """Shutdown normal servers"""
    num_vms = 10
    for i in range(1, num_vms + 1):
        print('Shutdown normal server for vm: ', i)
        os.system("ssh vm{} supervisorctl stop MP1_Test".format(i))
        os.system("ssh vm{} supervisorctl start MP1".format(i))


if __name__ == '__main__':
    num_vms = 10
    for i in range(1, num_vms + 1):
        print('Sending files for vm: ', i)
        vm_name = "vm{0}".format(str(i))
        send_files_to_vm(vm_name)

    shutdown_normal_servers_and_start_test_servers()

    # Known pattern on only one vm
    unit_test("qwertyvm2", '1')

    # Known pattern on two vms
    unit_test("zxcvvm2\|zxcvvm3", '2')

    # Known pattern on three vms
    unit_test("zxcvvm1\|zxcvvm2\|zxcvvm3", '3')

    # Known pattern on no vms
    unit_test("zxcvvm2\&zxcvvm3", '0')

    # Known pattern on all vms
    unit_test("qwertyvm", '10')

    # Known pattern on all vms (regex)
    unit_test("qwertyvm*", '10')

    shutdown_test_servers_and_start_normal_servers()
