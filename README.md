# Readme

Authored by Oliver Thomas, Carl-Emil Andersen and Rune Schuster

This is a project from the Course Distributed Systems and Security.

Overall Functionality:

- This project implements a Peer-to-Peer (P2P) network where peers can join, connect, and transact with each other.
- The network uses RSA encryption and signatures for secure communication.
- Each peer maintains a ledger containing account balances

Security:

- RSA encryption and signatures are used to secure communication and verify the authenticity of transactions.
- The software wallet requires a password to access the private key, making it more difficult for attackers to steal it.
- The password must be at least 12 characters long and complex to prevent brute-force attacks.
- The private key is encrypted on disk using the hashed password, adding another layer of protection.
