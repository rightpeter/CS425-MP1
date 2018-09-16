import os

def get_line_count(output):
    return 0

def unit_test(grep_args,expected_line_count):
    output = os.popen('./client.bin "{0}"'.format(grep_args)).read()
    line_count = get_line_count(output)
    if line_count!=expected_line_count:
        print("Unit test failed for arguments: {0}. Expected line count: {1} Got: {2}".format(grep_args, expected_line_count,line_count))
    print("Test passed!")

if __name__ == '__main__':  
    unit_test("vm2",2)
    unit_test("vm*",10)
    unit_test("vm",10)
