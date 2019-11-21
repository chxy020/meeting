package migrations

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Init_20191113_141139 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Init_20191113_141139{}
	m.Created = "20191113_141139"

	migration.Register("Init_20191113_141139", m)
}

// Run the migrations
func (m *Init_20191113_141139) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(DROP_TABLE)
	m.SQL(CREATE_TABLE)
}

// Reverse the migrations
func (m *Init_20191113_141139) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}

const (
	DROP_TABLE = `
        DROP TABLE IF EXISTS attribute_type;
        DROP TABLE IF EXISTS attribute;
        DROP TABLE IF EXISTS attendee;
    `

	CREATE_TABLE = `
        -- --------------------------------------------------
	    --  Table Structure for "seater/models.AttributeType"
	    -- --------------------------------------------------
	    CREATE TABLE IF NOT EXISTS attribute_type (
		    id int(11) NOT NULL AUTO_INCREMENT,
		    name varchar(255) NOT NULL DEFAULT '',
            is_hided boolean NOT NULL DEFAULT 0,
		    modifytime datetime NOT NULL,
            PRIMARY KEY (id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- ----------------------------------------------
	    --  Table Structure for "seater/models.Attribute"
	    -- ----------------------------------------------
	    CREATE TABLE IF NOT EXISTS attribute (
		    id int(11) NOT NULL AUTO_INCREMENT,
            member_id int(11) NOT NULL,
            attribute_type_id int(11) NOT NULL,
            parent_id int(11) DEFAULT NULL,
		    name varchar(127) NOT NULL DEFAULT '',
            content varchar(255) NOT NULL DEFAULT '',
            label text NOT NULL,
		    modifytime datetime NOT NULL,
            PRIMARY KEY (id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- ---------------------------------------------
	    --  Table Structure for "seater/models.Attendee"
	    -- ---------------------------------------------
	    CREATE TABLE IF NOT EXISTS attendee (
		    id int(11) NOT NULL AUTO_INCREMENT,
            meeting_id int(11) NOT NULL,
            rulesetup_id bigint(20) DEFAULT NULL,
            name varchar(255) NOT NULL DEFAULT '',
            company varchar(255) NOT NULL DEFAULT '',
            duties varchar(255) NOT NULL DEFAULT '',
            phone1 varchar(255) NOT NULL DEFAULT '',
            phone2 varchar(255) NOT NULL DEFAULT '',
            contacts varchar(255) NOT NULL DEFAULT '',
            contacts_phone varchar(255) NOT NULL DEFAULT '',
            card_id varchar(255) NOT NULL DEFAULT '',
            delstate int(11) DEFAULT NULL,
            image_url varchar(255) NOT NULL DEFAULT '',
            state int(11) DEFAULT NULL,
            compareimg1 varchar(255) NOT NULL DEFAULT '',
            compareimg2 varchar(255) NOT NULL DEFAULT '',
            compareimg3 varchar(255) NOT NULL DEFAULT '',
            camera varchar(255) NOT NULL DEFAULT '',
            seatid varchar(50) NOT NULL DEFAULT '',
            szm varchar(255) NOT NULL DEFAULT '',
            is_left int(11) DEFAULT NULL,
            xsorder int(11) DEFAULT NULL,
            viproom_id int(11) DEFAULT NULL,
            attributes text NOT NULL,
            modifytime datetime NOT NULL,
            PRIMARY KEY (id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `
)