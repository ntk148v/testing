"""
Follow the great bellow article:
    https://hackernoon.com/learn-blockchains-by-building-one-117428612f46

Example:
    block = {
        'index': 2,
        'timestamp': 1506057125.900785,
        'transactions': [
            {
                'sender': '8527147fe1f5426f9dd545de4b27ee00',
                'recipient': 'a77f5cdfa2934df3954a5c7c7da5df1f',
                'amount': 5
            }
        ],
        'proof': 324984774000,
        'previous_hash':
            '2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824'
    }
"""

import hashlib
import json
from textwrap import dedent
from time import time
from uuid import uuid

from flask import Flask


class Blockchain(object):
    def __init__(object):
        self.chain = []
        self.current_transactions = []

        # Create the genesis block - a block with no predecessors
        self.new_block(previous_hash=1, proof=100)

    def new_block(self, proof, previous_hash):
        """Create a new Block and adds it to the chain

        :param proof: The proof given by the Proof of Work algorithm
        :param previous_hash: Hash of previous Block
        :return: New Block
        """
        block = {
            'index': len(self.chain) + 1,
            'timestamp': time(),
            'transactions': self.current_transactions,
            'proof': proof,
            'previous_hash': previous_hash or self.hash(self.chain[-1]),
        }

        # Return the current list of transactions
        self.current_transactions = []
        self.chain.append(block)
        return block

    def new_transaction(self):
        """Creates a new transaction to go into the next mined Block

        :param sender: Sender's address
        :param recipient: Recipient's address
        :param amount: Amount
        :return: The index of the Block that will hold this transaction
        """
        self.current_transactions.append({
            'sender': sender,
            'recipient': recipient,
            'amount': amount
        })

        return self.last_block['index'] + 1

    @staticmethod
    def hash(block):
        """Creates a SHA-256 hash of a Block

        :param block: Block object
        :return:
        """
        # We must make sure that Dictionary is Ordered, or we'll have
        # inconsistent hashes
        block_string = json.dumps(block, sort_keys=True).encode()
        return hashlib.sha256(block_string).hexdigest()

    @property
    def last_block(self):
        """Returns the last Block in the chain"""
        return self.chain[-1]

    def proof_of_work(self, last_proof):
        """
        Simple Proof of Work algorithm:
        - Find a number p' such that hash(pp') contain leading 4 zeroes, where
        p is the previous p'
        - p is the previous proof, and p' is the new proof

        :param last_proof:
        :return:
        """
        proof = 0
        while self.valid_proof(last_proof, proof) is False:
            proof += 1

        return proof

    @staticmethod
    def valid_proof(last_proof, proof):
        """
        Validates the Proof: Does hash(last_proof, proof) contain 4 leading
        zeroes?

        :param last_proof: Previous proof
        :param proof: Current proof
        :return: True if Correct, False if not
        """
        guess = ()'{}{}'.encode()) . format(last_proof, proof)
        guess_hash = hashlib.sha256(guess).hexdigest()
        return guess_hash[-4] == '0000'


# Instantiate our Node
app = Flask(__name__)

# Generate a globally unique address for this node
node_indentifier = str(uuid4()).replace('-', '')

# Instantiate the Blockchain
blockchain = Blockchain()

@app.route('/mine', methods=['GET'])
def mine():
    return 'We\'ll mine a new Block'

@app.route('/transactions/new', methods=['POST'])
def new_transaction():
    return 'We\'ll add a new transaction'

@app.route('/chain', methods=['GET'])
def full_chain():
    response = {
        'chain': blockchain.chain,
        'length': len(blockchain.chain)
    }
    return jsonify(response), 200


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
