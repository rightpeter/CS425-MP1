import os
import string
import random

# reference: http://www.bswen.com/2018/04/python-How-to-generate-random-large-file-using-python.html
def generate_random_file(filename, size, know_pattern):
    random_chars = [random.choice(string.letters) for i in xrange(size)]
    random_ind= random.randint(0,size)
    # Insert new lines randomly
    i=0
    while i<size:
        i+=random.randint(10,30)
        random_chars.insert(i,'\n')
    # Insert known pattern
    random_chars.insert(random_ind,'\n'+know_pattern+'\n')
    with open(filename, 'w+') as f:
        f.write(''.join(random_chars))

def send_log_file(vm_name, file_name):
    os.system("scp ./{0} {1}:/tmp/".format(file_name, vm_name))

def send_config_file(vm):
    os.system("scp ./config.json {1}:/tmp/".format(vm))

def send_files_to_vm(vm):
    log_file_name="{0}.log".format(vm)
    known_pattern = "qwerty{0}qwerty\n{0}".format(vm)
    generate_random_file(log_file_name,1024,known_pattern)
    send_log_file(vm,log_file_name)
    send_config_file(vm)

if __name__ == '__main__':
    
    num_vms = 10
    for i in xrange(1,num_vms+1):
        vm_name = "vm{0}".format(str(i))
        send_files_to_vm(vm_name)

# scp ./test.txt vm1:/tmp