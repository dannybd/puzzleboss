#!/usr/bin/python
from __future__ import print_function
import httplib2
import os

import googleapiclient 
from googleapiclient import discovery
import oauth2client
from oauth2client import file
from oauth2client import client
from oauth2client import tools

flags = None

SCOPES = 'https://www.googleapis.com/auth/admin.directory.user'
CLIENT_SECRET_FILE = 'client_secret.json'
APPLICATION_NAME = 'Puzzlebitch Scripts'


def get_credentials():
    """Gets valid user credentials from storage.

    If nothing has been stored, or if the stored credentials are invalid,
    the OAuth2 flow is completed to obtain the new credentials.

    Returns:
        Credentials, the obtained credential.
    """
    credential_dir = '/canadia/puzzlebitch/googcredentials'
    if not os.path.exists(credential_dir):
        os.makedirs(credential_dir)
    credential_path = os.path.join(credential_dir,
                                   'admin-directory_v1-python.json')

    store = oauth2client.file.Storage(credential_path)
    credentials = store.get()
    if not credentials or credentials.invalid:
        flow = client.flow_from_clientsecrets(CLIENT_SECRET_FILE, SCOPES)
        flow.user_agent = APPLICATION_NAME
        credentials = oauth2client.tools.run_flow(flow, store, flags)
        print('Storing credentials to ' + credential_path)
    return credentials

def main():
    """
    Creates a Google Admin SDK API service object and outputs a list of first
    200 users in the domain.
    """
    credentials = get_credentials()
    http = credentials.authorize(httplib2.Http())
    service = discovery.build('admin', 'directory_v1', http=http)

    print('Getting the first 200 users in the domain')
    results = service.users().list(customer='my_customer', maxResults=200,
        orderBy='email').execute()
    users = results.get('users', [])

    if not users:
        print('No users in the domain.')
    else:
        print('Users:')
        for user in users:
            print('{0} ({1})'.format(user['primaryEmail'],
                user['name']['fullName']))


if __name__ == '__main__':
    main()

