/**
 * CS61065 - Assignment 4 - Part C
 * 
 * Authors:
 * Utkarsh Patel (18EC35034)
 * Saransh Patel (18CS30039)
 */

const { Gateway, Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');

const fs = require('fs');
const path = require('path');
const prompt = require('prompt-sync')({sigint: true});

async function main() {
    /* Org1 connection profile */
    console.log('[*] Creating Org1 connection profile..');
    const ccpPath1 = path.resolve('../organizations/peerOrganizations/org1.example.com/connection-org1.json');
    const ccp1 = JSON.parse(fs.readFileSync(ccpPath1, 'utf8'));
    console.log('[+] Created Org1 connection profile');


    /* Org2 connection profile */
    console.log('[*] Creating Org2 connection profile..');
    const ccpPath2 = path.resolve('../organizations/peerOrganizations/org2.example.com/connection-org2.json');
    const ccp2 = JSON.parse(fs.readFileSync(ccpPath2, 'utf8'));
    console.log('[+] Created Org2 connection profile');

    /* Org1 certificate authorities */
    console.log('[*] Creating Org1 certificate authorities..')
    const caInfo1 = ccp1.certificateAuthorities['ca.org1.example.com'];
    const caTLSCACerts1 = caInfo1.tlsCACerts.pem;
    const ca1 = new FabricCAServices(
        caInfo1.url, 
        {
            trustedRoots: caTLSCACerts1,
            verify: false
        },
        caInfo1.caName
    );
    console.log('[+] Created Org1 certificate authorities')

    /* Org2 certificate authorities */
    console.log('[*] Creating Org2 certificate authorities..')
    const caInfo2 = ccp2.certificateAuthorities['ca.org2.example.com'];
    const caTLSCACerts2 = caInfo2.tlsCACerts.pem;
    const ca2 = new FabricCAServices(
        caInfo2.url, 
        {
            trustedRoots: caTLSCACerts2,
            verify: false
        },
        caInfo2.caName
    );
    console.log('[+] Created Org2 certificate authorities')

    /* Creating wallet for Org1*/
    console.log('[*] Creating wallet for Org1..')
    const walletPath1 = path.join(process.cwd(), 'wallet1');
    const wallet1 = await Wallets.newFileSystemWallet(walletPath1);
    console.log('[+] Created wallet for Org1')

    /* Creating wallet for Org2*/
    console.log('[*] Creating wallet for Org2..')
    const walletPath2 = path.join(process.cwd(), 'wallet2');
    const wallet2 = await Wallets.newFileSystemWallet(walletPath2);
    console.log('[+] Created wallet for Org2')

    /* Get admin identity  for Org1 */
    console.log('[*] Creating admin for Org1..')
    var adminID1 = await wallet1.get("admin");
    const enrollment1 = await ca1.enroll({
        enrollmentID: 'admin',
        enrollmentSecret: 'adminpw'
    });
    const x509ID1 = {
        credentials: {
            certificate: enrollment1.certificate,
            privateKey: enrollment1.key.toBytes()
        },
        mspId: 'Org1MSP',
        type: 'X.509'
    };
    await wallet1.put("admin", x509ID1);
    adminID1 = await wallet1.get("admin");
    console.log('[+] Created admin for Org1')

    /* Get admin identity for Org2 */
    console.log('[*] Creating admin for Org2..')
    var adminID2 = await wallet2.get("admin");
    const enrollment2 = await ca2.enroll({
        enrollmentID: 'admin',
        enrollmentSecret: 'adminpw'
    });
    const x509ID2 = {
        credentials: {
            certificate: enrollment2.certificate,
            privateKey: enrollment2.key.toBytes()
        },
        mspId: 'Org2MSP',
        type: 'X.509'
    };
    await wallet2.put("admin", x509ID2);
    adminID2 = await wallet2.get("admin");
    console.log('[+] Created admin for Org2')

    /* User registration for Org1 */
    var userID1 = await wallet1.get("appUser");
    if (!userID1) {
        console.log('[*] Registering peer0 from Org1..')
        const provider1 = wallet1.getProviderRegistry().getProvider(adminID1.type);
        const adminUser1 = await provider1.getUserContext(adminID1, "admin");

        const secret1 = await ca1.register({
            affiliation: 'org1.department1',
            enrollmentID: 'appUser',
            role: 'client'
        }, adminUser1);

        const enrollment1 = await ca1.enroll({
            enrollmentID: 'appUser',
            enrollmentSecret: secret1
        });

        const x509ID1 = {
            credentials: {
                certificate: enrollment1.certificate,
                privateKey: enrollment1.key.toBytes()
            },
            mspId: 'Org1MSP',
            type: 'X.509'
        };

        await wallet1.put('appUser', x509ID1);
        userID1 = await wallet1.get('appUser');
        console.log('[+] peer0 from Org1 registered')
    }

    /* User registration for Org2 */
    var userID2 = await wallet2.get("appUser");
    if (!userID2) {
        console.log('[*] Registering peer0 from Org2..')
        const provider2 = wallet2.getProviderRegistry().getProvider(adminID2.type);
        const adminUser2 = await provider2.getUserContext(adminID2, "admin");

        const secret2 = await ca2.register({
            affiliation: 'org2.department1',
            enrollmentID: 'appUser',
            role: 'client'
        }, adminUser2);

        const enrollment2 = await ca2.enroll({
            enrollmentID: 'appUser',
            enrollmentSecret: secret2
        });

        const x509ID2 = {
            credentials: {
                certificate: enrollment2.certificate,
                privateKey: enrollment2.key.toBytes()
            },
            mspId: 'Org2MSP',
            type: 'X.509'
        };

        await wallet2.put('appUser', x509ID2);
        userID2 = await wallet2.get('appUser');
        console.log('[+] peer0 from Org2 registered')
    }

    /* Connect to gateway for Org1 */
    console.log('[*] Connecting to gateway for Org1')
    const gateway1 = new Gateway();
    await gateway1.connect(
        ccp1,
        {
            wallet: wallet1,
            identity: 'appUser',
            discovery: {
                enabled: true,
                asLocalhost: true
            }
        }
    );
    console.log('[+] Connected to gateway for Org1')

    /* Connect to gateway for Org2*/
    console.log('[*] Connecting to gateway for Org2')
    const gateway2 = new Gateway();
    await gateway2.connect(
        ccp2,
        {
            wallet: wallet2,
            identity: 'appUser',
            discovery: {
                enabled: true,
                asLocalhost: true
            }
        }
    );
    console.log('[+] Connected to gateway for Org2')

    /* Connect to channel */
    const network1 = await gateway1.getNetwork('mychannel');
    const network2 = await gateway2.getNetwork('mychannel');

    /* Specifying chaincode */
    const contract1 = network1.getContract('bstchaincode');
    const contract2 = network2.getContract('bstchaincode');

    let turn = 1, val, result, flag = true;
    const contracts = [contract1, contract2];

    while (flag) {
        turn ^= 1   /* Switch peer */

        const cmd = prompt(`[peer0.org${turn + 1}] Enter command: `);
        switch (cmd) {
            case 'INSERT':
                val = prompt(`[peer0.org${turn + 1}] Enter value to insert: `);
                val = Number(val);
                try {
                    await contracts[turn].submitTransaction('Insert', val);
                    console.log(`[+] Insertion completed successfully.`)
                } catch (err) {
                    console.log(`[-] Error: ${err.message}`)
                }
                break;
            
            case 'DELETE':
                val = prompt(`[peer0.org${turn + 1}] Enter value to delete: `);
                val = Number(val);
                try {
                    await contracts[turn].submitTransaction('Delete', val);
                    console.log(`[+] Deletion completed successfully.`)
                } catch (err) {
                    console.log(`[-] Error: ${err.message}`)
                }
                break;

            case 'INORDER':
                try {
                    result = await contracts[turn].evaluateTransaction('Inorder');
                    console.log(`[+] Inorder traversal: ${result.toString()}`);
                } catch (err) {
                    console.log(`[-] Error: ${err.message}`)
                }
                break;

            case 'PREORDER':
                try {
                    result = await contracts[turn].evaluateTransaction('Preorder');
                    console.log(`[+] Preorder traversal: ${result.toString()}`);
                } catch (err) {
                    console.log(`[-] Error: ${err.message}`)
                }
                break;

            case 'TREEHEIGHT':
                try {
                    result = await contracts[turn].evaluateTransaction('TreeHeight');
                    console.log(`[+] Tree height: ${result.toString()}`);
                } catch (err) {
                    console.log(`[-] Error: ${err.message}`);
                }
                break;

            case 'EXIT':
                console.log('[+] Closing..');
                flag = false;
                break;

            default:
                console.log('[-] Invalid command!\n');
                turn ^= 1
        }
    }

    gateway1.disconnect();
    gateway2.disconnect();
}

main();