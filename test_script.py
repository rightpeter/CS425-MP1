import os

def get_line_count(output):
    return int(output.split('Total line count:')[-1])

def unit_test(grep_args,expected_line_count):
    output = os.popen('./client.bin "{0}"'.format(grep_args)).read()
    line_count = get_line_count(output)
    print("Testing pattern: {0}".format(grep_args))
    if line_count!=expected_line_count:
        print("Unit test failed for arguments: {0}. Expected line count: {1} Got: {2}".format(grep_args, expected_line_count,line_count))
        return
    print("Test passed!")

if __name__ == '__main__':  

    # Known pattern on only one vm
    unit_test("qwertyvm2",1)

    # Known pattern on two vms
    unit_test("zxcvvm2\|zxcvvm3",2)

    # Known pattern on three vms
    unit_test("zxcvvm1\|zxcvvm2\|zxcvvm3",3)

    # Known pattern on no vms
    unit_test("zxcvvm2\&zxcvvm3",0)

    # Known pattern on all vms
    unit_test("quertyvm",10)

    # Known pattern on all vms (regex)
    unit_test("quertyvm*",10)


