#!/bin/bash

export PERL5LIB='.:/canadia/puzzleboss/lib'

eval $(perl -MPB::Config -e 'PB::Config::export_to_bash();')

deleteuser=`echo $1 | perl -pi -e 's/[[:space:]*]//g'`


if [ -z "$deleteuser" ]; then
    echo "please specify user to delete"
else

#    echo "removing user $deleteuser from HuntTeam LDAP group"
#    /canadia/puzzleboss/admin-tools/ldapdeleteuserfromgroup $deleteuser HuntTeam

    echo "deleting user $deleteuser"
    deletedn=`ldapsearch -h ${REGISTER_LDAP_HOST} -xLLL -b "${REGISTER_LDAP_DC}" uid=$deleteuser dn|cut -d':' -f2|head -n 1|perl -pi -e 's/[[:space:]*]//g'`
    if [ -z "$deletedn" ]; then
	echo "ldap user does not exist"
    else 
	echo "deleting ldap dn $deletedn"
	ldapdelete -h ${REGISTER_LDAP_HOST} -x -D cn=${REGISTER_LDAP_ADMIN_USER},${REGISTER_LDAP_DC} -w ${REGISTER_LDAP_ADMIN_PASS} $deletedn
    fi


    echo "deleting google user $deleteuser"
    echo "with google domain ${GOOGLE_DOMAIN}"
    (cd $PB_PATH/newgoogle && ./DeleteDomainUser.py ${deleteuser}@${GOOGLE_DOMAIN} )


#    twikitopic=/canadia/twiki/data/Main/$deleteuser.txt
#    echo "deleting twiki topic file $twikitopic"
#    rm -rf $twikitopic
#    rm -rf $twikitopic,v

fi

exit 0
