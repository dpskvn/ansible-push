# ansible-push
This basic connector implements `testconnection` and `installcertificatebundle` routes. Both communicate with AAP/AWX, and the installation one is capable of starting a job template with a specific ID.

Certificates are passed as extra variables to the job. As AAP doesn't support arrays as parameters, the certificate chain is passed as a `|`-delimited string.

## Sample Ansible playbook for testing
**Note:** Survery variables need to be configured on the AAP/AWX side.
`
```yaml
- name: Extract and display certificate information
  hosts: localhost
  gather_facts: false
  tasks:
    - name: Decode base64-encoded certificate
      ansible.builtin.copy:
        dest: "/tmp/temp_cert.pem"
        content: "{{ certificate | b64decode }}"
      register: cert_file

    - name: Retrieve certificate information
      community.crypto.x509_certificate_info:
        path: "/tmp/temp_cert.pem"
      register: cert_info

    - name: Display certificate information
      ansible.builtin.debug:
        var: cert_info
```
## TODO
- [x] Implement a basic PoC
- [x] Create a basic `manifest.json`
- [x] Clean up unnecessary types
- [x] Figure out how to handle `installcertificatebundle` vs `configureinstallationendpoint`
- [x] Remove certificate discovery as it doesn't make sense Ansible-wise
- [x] Finalize the parameters in `manifest.json`
