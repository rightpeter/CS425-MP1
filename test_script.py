import os


def unit_test(grep_args,expected_output):
    output = os.popen('./client.bin "{0}"'.format(grep_args)).read()
    if output!=expected_output:
        print("Unit test failed for arguments: {0}. Expected : {1} Got output: {2}".format(grep_args,expected_output,output))


if __name__ == '__main__':
    unit_test("MP1","MP1\nMP1")