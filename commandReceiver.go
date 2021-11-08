//Command receiver: takes newline split input from standard input and returns command as string
//	^^^ this will be replaced with tcp protocol that receives commands from client's command line
//Implement commands: Takes string and decodes it and finds out what commands to run and runs them
//if commands were able to be implemented, returns RESP okay to stdout
//if functions return errors, returns error message to RESP