# Devnet
Devnet was settup for benchmark, so the way to deploy is totally different with Testnet. 
We used Ansible to handle bunch of nodes (clear data, start, stop, ...) and it isn't easy for setup. For this reason; I only guide you how to build, push docker image & use Ansible to benchmark.

## 1.Building image
**Preparing:**
- Because this project is private, so we need token to clone it for building.
- Change the `login` & `password` (is your token) in `dockerfiles/node/token` to yours.  
- [Here](https://github.com/settings/tokens) setup your token.    


**Building**
- Run file `deploy/devnet/node/build_image.sh` it will clone latest code at develop branch & build it to image `kybernetwork/evrynet-node:1.0.1-dev`  
- After completing, you can use command `docker images -a` to check.
- Use `docker push kybernetwork/evrynet-node:1.0.1-dev` to push this image to docker hub.
- Nodes will pull this image to run.

## 2.Preparing Ansible to manage nodes
- Install Ansible on local [Here](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)
- Clone https://github.com/KyberNetwork/ansible-playbook/tree/evrynet on `evrynet` branch
- Create new file at `~/.ansible/vault-evrynet` with a content is `613bba68196d6937c7b31284377b3a959f5618004ddbfe5bccf7fcf12b01a3e8`

## 3.Using Ansible
- To reset data of a specific node, change field `evrynet_node_reset_db` to `true` at [this file](https://github.com/KyberNetwork/ansible-playbook/blob/evrynet/roles/evrynet-node-lite/defaults/main.yml#L16). Then run `ansible-playbook --vault-id ~/.ansible/vault-evrynet -i inventories/evrynet/hosts.yml evrynet-lite.yml -l dev-evrynet-node-01` 
- This command also pull new docker image from docker hub.
- To reset data of all node, remove `-l dev-evrynet-node-01`. It will clear data each node sequentially.
- After clearing data, you must revert `evrynet_node_reset_db` to `false` & run `ansible-playbook --vault-id ~/.ansible/vault-evrynet -i inventories/evrynet/hosts.yml evrynet-lite.yml` to restart nodes.
- If you don't revert `evrynet_node_reset_db` to `false`, nodes won't be started.
- To change the command to run node, you can edit at [this file](https://github.com/KyberNetwork/ansible-playbook/blob/evrynet/roles/evrynet-node-lite/templates/evrynet-node.service.j2)