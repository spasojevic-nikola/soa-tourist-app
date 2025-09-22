// sluzi za mapiranje korisnika 
// kad se registruje da se upise i u stakeholder users tabelu
package main

type StakeholderUser struct {
    ID           uint   `json:"id"`
    FirstName    string `json:"first_name"`
    LastName     string `json:"last_name"`
    Username     string `json:"username"`
    Role         string `json:"role"`
    ProfileImage string `json:"profile_image"`
    Biography    string `json:"biography"`
    Motto        string `json:"motto"`
}