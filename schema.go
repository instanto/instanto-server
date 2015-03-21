package main

import (
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func createSchema() {
	schema := `SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;


CREATE TABLE IF NOT EXISTS article (
id int(11) NOT NULL,
  title varchar(200) NOT NULL,
  web varchar(200) NOT NULL,
  date int(11) NOT NULL,
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL,
  newspaper int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS category (
id int(11) NOT NULL,
  name varchar(45) NOT NULL,
  description varchar(45) NOT NULL,
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS category_resource (
  category int(11) NOT NULL,
  resource varchar(45) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS config (
  config_key varchar(45) NOT NULL,
  config_value varchar(200) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS financed_project (
id int(11) NOT NULL,
  title varchar(45) NOT NULL,
  started int(11) NOT NULL,
  ended int(11) NOT NULL,
  budget int(11) NOT NULL,
  scope varchar(45) NOT NULL,
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL,
  primary_funding_body int(11) NOT NULL,
  primary_record varchar(45) NOT NULL,
  primary_leader int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS financed_project_leader (
  financed_project int(11) NOT NULL,
  member int(11) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS financed_project_member (
  financed_project int(11) NOT NULL,
  member int(11) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS funding_body (
id int(11) NOT NULL,
  name varchar(200) NOT NULL,
  web varchar(200) NOT NULL,
  scope varchar(45) NOT NULL,
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS funding_body_financed_project (
  funding_body int(11) NOT NULL,
  financed_project int(11) NOT NULL,
  record varchar(45) NOT NULL,
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS member (
id int(11) NOT NULL,
  first_name varchar(45) NOT NULL,
  last_name varchar(45) NOT NULL,
  degree varchar(45) CHARACTER SET big5 NOT NULL,
  year_in int(11) NOT NULL,
  year_out int(11) NOT NULL,
  email varchar(45) NOT NULL,
  cv varchar(32) NOT NULL DEFAULT '',
  photo varchar(32) NOT NULL DEFAULT '',
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL,
  primary_status int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS member_publication (
  member int(11) NOT NULL,
  publication int(11) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS member_status (
  member int(11) NOT NULL,
  status int(11) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS newspaper (
id int(11) NOT NULL,
  name varchar(200) NOT NULL,
  web varchar(200) NOT NULL,
  logo varchar(32) NOT NULL DEFAULT '',
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS partner (
id int(11) NOT NULL,
  name varchar(45) NOT NULL,
  web varchar(45) NOT NULL,
  logo varchar(32) NOT NULL DEFAULT '',
  same_department tinyint(1) NOT NULL,
  scope varchar(45) NOT NULL,
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='z' AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS partner_member (
  partner int(11) NOT NULL,
  member int(11) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS permission (
  id varchar(64) NOT NULL,
  display_name varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO permission (id, display_name) VALUES
('access_log_view', 'View the access log'),
('article_create', 'Create article'),
('article_delete_mine', 'Delete article created by me'),
('article_delete_other', 'Delete article created by other'),
('article_update_mine', 'Update article created by me'),
('article_update_other', 'Update article created by other'),
('financed_project_create', 'Create financed project'),
('financed_project_delete_mine', 'Delete financed project created by me'),
('financed_project_delete_other', 'Delete financed project created by other'),
('financed_project_leader_add_mine_mine', 'Add leaders created by me to a financed project created by me'),
('financed_project_leader_add_mine_other', 'Add leaders created by other to a financed project created by me'),
('financed_project_leader_add_other_mine', 'Add leaders created by me to a financed project created by other'),
('financed_project_leader_add_other_other', 'Add leaders created by other to a financed project created by other'),
('financed_project_leader_remove_mine_mine', 'Remove leaders created by me from a financed project created by me'),
('financed_project_leader_remove_mine_other', 'Remove members created by other from a financed project created by me'),
('financed_project_leader_remove_other_mine', 'Remove members created by other from a financed project created by me'),
('financed_project_leader_remove_other_other', 'Remove members created by other from a financed project created by other'),
('financed_project_member_add_mine_mine', 'Add members created by me to a financed project created by me'),
('financed_project_member_add_mine_other', 'Add members created by other to a financed project created by other'),
('financed_project_member_add_other_mine', 'Add members created by me to a financed project created by other'),
('financed_project_member_add_other_other', 'Add members created by other to a financed project created by other'),
('financed_project_member_remove_mine_mine', 'Remove members created by me from a financed project created by me'),
('financed_project_member_remove_mine_other', 'Remove members created by other from a financed project created by me'),
('financed_project_member_remove_other_mine', 'Remove members created by me from a financed project created by other'),
('financed_project_member_remove_other_other', 'Remove members created by other from a financed project created by other'),
('financed_project_update_mine', 'Update financed project created by me'),
('financed_project_update_other', 'Update financed project created by other'),
('funding_body_create', 'Create funding body'),
('funding_body_delete_mine', 'Delete funding body created by me'),
('funding_body_delete_other', 'Delete funding body created by other'),
('funding_body_financed_project_add_mine_mine', 'Add financed projects created by me to a funding body created by me'),
('funding_body_financed_project_mine_other', 'Add financed projects created by other to a funding body created by me'),
('funding_body_financed_project_other_mine', 'Add financed projects created by me to a funding body created by other'),
('funding_body_financed_project_other_other', 'Add financed projects created by other to a funding body created by other'),
('funding_body_financed_project_remove_mine_mine', 'Remove financed projects created by me from a funding body created by me'),
('funding_body_financed_project_remove_mine_other', 'Remove financed projects created by other from a funding body created by me'),
('funding_body_financed_project_remove_other_mine', 'Remove financed projects created by me from a funding body created by other'),
('funding_body_financed_project_remove_other_other', 'Remove financed projects created by other from a funding body created by other'),
('funding_body_financed_project_update_mine_mine', 'Update the record of the funding if the funding body and the financed project are mine'),
('funding_body_update_mine', 'Update funding body created by me'),
('funding_body_update_other', 'Update funding body created by other'),
('group_create', 'Create group'),
('group_delete_mine', 'Delete group created by me'),
('group_delete_other', 'Delete group created by other'),
('group_update_mine', 'Update group created by me'),
('group_update_other', 'Update group created by other'),
('group_user_add_mine_mine', 'Add users created by me to a group created by me'),
('group_user_add_mine_other', 'Add users created by other to a group created by me'),
('group_user_add_other_mine', 'Add users created by me to a group created by other'),
('group_user_add_other_other', 'Add users created by other to a group created by other'),
('group_user_remove_mine_mine', 'Remove users created by me from a group created by me'),
('group_user_remove_mine_other', 'Remove users created by other from a group created by me'),
('group_user_remove_other_mine', 'Remove users created by other from a group created by me'),
('group_user_remove_other_other', 'Remove users created by other from a group created by other'),
('log_view', 'View the system log'),
('member_create', 'Create member'),
('member_delete_mine', 'Delete member created by me'),
('member_delete_other', 'Delete member created by other'),
('member_publication_add_mine_mine', 'Add publications created by me to a member created by me'),
('member_publication_add_mine_other', 'Add publications created by other to a member created by me'),
('member_publication_add_other_mine', 'Add publications created by me to a member created by other'),
('member_publication_add_other_other', 'Add publications created by other to a member created by other'),
('member_publication_remove_mine_mine', 'Remove publications created by me from a member created by me'),
('member_publication_remove_mine_other', 'Remove publications created by other from a member created by me'),
('member_publication_remove_other_mine', 'Remove publications created by me from a member created by other'),
('member_publication_remove_other_other', 'Remove publications created by other from a member created by other'),
('member_status_add_mine_mine', 'Add status created by me to a member created by me'),
('member_status_add_mine_other', 'Add statuses created by other to a member created by me'),
('member_status_add_other_mine', 'Add statuses created by me to a member created by other'),
('member_status_add_other_other', 'Add statuses created by other to a member created by other'),
('member_status_remove_mine_mine', 'Remove status created by me from a member created by me'),
('member_status_remove_mine_other', 'Remove statuses created by other from a member created by me'),
('member_status_remove_other_mine', 'Remove statuses created by me from a member created by other'),
('member_status_remove_other_other', 'Remove statuses created by other from a member created by other'),
('member_update_mine', 'Update member created by me'),
('member_update_other', 'Update member created by other'),
('newspaper_create', 'Create newspaper'),
('newspaper_delete_mine', 'Delete newspaper created by me'),
('newspaper_delete_other', 'Delete newspaper created by other'),
('newspaper_update_mine', 'Update newspaper created by me'),
('newspaper_update_other', 'Update newspaper created by other'),
('partner_create', 'Create partner'),
('partner_delete_mine', 'Delete partner created by me'),
('partner_delete_other', 'Delete newspaper created by other'),
('partner_member_add_mine_mine', 'Add members created by me to a partner created by me'),
('partner_member_add_mine_other', 'Add members created by other to a partner created by me'),
('partner_member_add_other_mine', 'Add members created by me to a partner created by other'),
('partner_member_add_other_other', 'Add members created by other to a partner created by other'),
('partner_member_remove_mine_mine', 'Remove members created by me from a partner created by me'),
('partner_member_remove_mine_other', 'Remove members created by other from a partner created by me'),
('partner_member_remove_other_mine', 'Remove members created by me from a partner created by other'),
('partner_member_remove_other_other', 'Remove members created by other from a partner created by other'),
('partner_update_mine', 'Update partner created by me'),
('partner_update_other', 'Update newspaper created by other'),
('publication_create', 'Create publication'),
('publication_delete_mine', 'Delete publication created by me'),
('publication_delete_other', 'Delete publication created by other'),
('publication_type_create', 'Create publication type'),
('publication_type_delete_mine', 'Delete publication type created by me'),
('publication_type_delete_other', 'Delete publication created by other'),
('publication_type_update_mine', 'Update publication type created by me'),
('publication_type_update_other', 'Update publication created by other'),
('publication_update_mine', 'Update publication created by me'),
('publication_update_other', 'Update publication created by other'),
('publisher_create', 'Create publisher'),
('publisher_delete_mine', 'Delete publisher created by me'),
('publisher_delete_other', 'Delete publisher created by other'),
('publisher_update_mine', 'Update publisher created by me'),
('publisher_update_other', 'Update publisher created by other'),
('research_area_create', 'Create research area'),
('research_area_delete_mine', 'Delete research area created by me'),
('research_area_delete_other', 'Delete research area created by other'),
('research_area_research_line_add_mine_mine', 'Add research lines created by me to a research area created by me'),
('research_area_research_line_add_mine_other', 'Add research lines created by other to a research area created by me'),
('research_area_research_line_add_other_mine', 'Add research lines created by me to a research area created by other'),
('research_area_research_line_add_other_other', 'Add research lines created by other to a research area created by other'),
('research_area_research_line_remove_mine_mine', 'Remove research lines created by me from a research area created by me'),
('research_area_research_line_remove_mine_other', 'Remove research lines created by other from a research area created by me'),
('research_area_research_line_remove_other_mine', 'Remove research lines created by me from a research area created by other'),
('research_area_research_line_remove_other_other', 'Remove research lines created by other from a research area created by other'),
('research_area_update_mine', 'Update research area created by me'),
('research_area_update_other', 'Update research area created by other'),
('research_line_article_add_mine_mine', 'Add articles created by me to a research line created by me'),
('research_line_article_add_mine_other', 'Add articles created by other to a research line created by me'),
('research_line_article_add_other_mine', 'Add articles created by me to a research linecreated by other'),
('research_line_article_add_other_other', 'Add articles created by other to a research line created by other'),
('research_line_article_remove_mine_mine', 'Remove articles created by me from a research line created by me'),
('research_line_article_remove_mine_other', 'Remove articles created by other from a research line created by me'),
('research_line_article_remove_other_mine', 'Remove articles created by me from a research line created by other'),
('research_line_article_remove_other_other', 'Remove articles created by other from a research line created by other'),
('research_line_create', 'Create research line'),
('research_line_delete_mine', 'Delete research line created by me'),
('research_line_delete_other', 'Delete research line created by other'),
('research_line_financed_project_add_mine_mine', 'Add financed projects created by me to a research line created by me'),
('research_line_financed_project_add_mine_other', 'Add financed projects created by other to a research line created by me'),
('research_line_financed_project_add_other_mine', 'Add financed projects created by me to a research line created by other'),
('research_line_financed_project_add_other_other', 'Add financed projects created by other to a research line created by other'),
('research_line_financed_project_remove_mine_mine', 'Remove financed projects created by me from a research line created by me'),
('research_line_financed_project_remove_mine_other', 'Remove financed projects created by other from a research line created by me'),
('research_line_financed_project_remove_other_mine', 'Remove financed projects created by me from a research line created by other'),
('research_line_financed_project_remove_other_other', 'Remove financed projects created by other from a research line created by other'),
('research_line_member_add_mine_mine', 'Add members created by me to a research line created by me'),
('research_line_member_add_mine_other', 'Add members created by other to a research line created by me'),
('research_line_member_add_other_mine', 'Add members created by me to a research line created by other'),
('research_line_member_add_other_other', 'Add members created by other to a research line created by other'),
('research_line_member_remove_mine_mine', 'Remove members created by me from a research line created by me'),
('research_line_member_remove_mine_other', 'Remove members created by other from a research line created by me'),
('research_line_member_remove_other_mien', 'Remove members created by me from a research line created by other'),
('research_line_member_remove_other_other', 'Remove members created by other from a research line created by other'),
('research_line_partner_add_mine_mine', 'Add partners created by me to a research line created by me'),
('research_line_partner_add_mine_other', 'Add partners created by other from a research line created by me'),
('research_line_partner_add_other_mine', 'Add partners created by me from a research line created by other'),
('research_line_partner_add_other_other', 'Add partners created by other from a research line created by other'),
('research_line_partner_remove_mine_mine', 'Remove partners created by me from a research line created by me'),
('research_line_partner_remove_mine_other', 'Remove partners created by other from a research line created by me'),
('research_line_partner_remove_other_mine', 'Remove partners created by me from a research line created by other'),
('research_line_partner_remove_other_other', 'Remove partners created by other from a research line created by other'),
('research_line_publication_add_mine_mine', 'Add publications created by me to a research line created by me'),
('research_line_publication_add_mine_other', 'Add publications created by other to a research line created by me'),
('research_line_publication_add_other_mine', 'Add publications created by me to a research line created by other'),
('research_line_publication_add_other_other', 'Add publication created by other to a research line created by other'),
('research_line_publication_remove_mine_mine', 'Remove publications created by me from a research line created by me'),
('research_line_publication_remove_mine_other', 'Remove publications created by other from a research line created by me'),
('research_line_publication_remove_other_mine', 'Remove publications created by me from a research line created by other'),
('research_line_publication_remove_other_other', 'Remove publications created by other from a research line creared by other'),
('research_line_resouce_remove_mine_other', 'Remove resources created by other from a research line created by me'),
('research_line_resource_add_mine_mine', 'Add resources created by me to a research line created by me'),
('research_line_resource_add_mine_other', 'Add resources created by other to a research line created by me'),
('research_line_resource_add_other_mine', 'Add resources created by me to a research line created by other'),
('research_line_resource_add_other_other', 'Add resources created by other to a reseach line created by other'),
('research_line_resource_remove_mine_mine', 'Remove resources created by me from a research line created by me'),
('research_line_resource_remove_other_mine', 'Remove resources created by me from a research line created by other'),
('research_line_resource_remove_other_other', 'Remove resources created by other from a research line created by other'),
('research_line_student_work_add_mine_mine', 'Add student works created by me to a research line created by me'),
('research_line_student_work_add_mine_other', 'Add student works created by other to a reseach line created by me'),
('research_line_student_work_add_other_mine', 'Add student works created by me to a research line created by other'),
('research_line_student_work_add_other_other', 'Add student works created by other to a research line created by other'),
('research_line_student_work_remove_mine_mine', 'Remove student workscreated by me  from a research line created by me'),
('research_line_student_work_remove_mine_other', 'Remove student works created by other from a research line created by me'),
('research_line_student_work_remove_other_mine', 'Remove student works created by me from a research line created by other'),
('research_line_student_work_remove_other_other', 'Remove student works created by other from a research line created by other'),
('research_line_update_mine', 'Update research line created by me'),
('research_line_update_other', 'Update research line created by other'),
('resource_create', 'Create reosource'),
('resource_delete_mine', 'Delete resource created by me'),
('resource_delete_other', 'Delete resource created by other'),
('resource_type_create', 'Create resource type'),
('resource_type_delete', 'Delete resource type created by other'),
('resource_type_delete_mine', 'Delete resource type created by me'),
('resource_type_update', 'Update resource type created by other'),
('resource_type_update_mine', 'Update resource type created by me'),
('resource_update_mine', 'Update resource created by me'),
('resource_update_other', 'Update resource created by other'),
('rol_create', 'Create rol'),
('rol_delete_mine', 'Delete rol created by me'),
('rol_delete_other', 'Delete rol created by other'),
('rol_permission_add_mine', 'Add permissions to a rol created by me'),
('rol_permission_add_other', 'Add permissions to a rol created by other'),
('rol_permission_remove_mine', 'Remove permissions from a rol created by me'),
('rol_permission_remove_other', 'Remove permissions from a rol created by other'),
('rol_update_mine', 'Update rol created by me'),
('rol_update_other', 'Update rol created by other'),
('status_create', 'Create status'),
('status_delete_mine', 'Delete status created by me'),
('status_delete_other', 'Delete status created by other'),
('status_update_mine', 'Update status created by me'),
('status_update_other', 'Update status created by other'),
('student_work_create', 'Create student work'),
('student_work_delete_mine', 'Delete student work created by me'),
('student_work_delete_other', 'Delete student work created by other'),
('student_work_type_create', 'Create student work type'),
('student_work_type_delete_mine', 'Delete student work type created by me'),
('student_work_type_delete_other', 'Delete student work type created by other'),
('student_work_type_update_mine', 'Update student work type created by me'),
('student_work_type_update_other', 'Update student work type created by other'),
('student_work_update_mine', 'Update student work created by me'),
('student_work_update_other', 'Update student work created by other'),
('user_create', 'Create user'),
('user_disable_itself', 'Disable the user logged in only'),
('user_disable_mine', 'Disable user created by me'),
('user_disable_other', 'Disable user created by other'),
('user_enable_mine', 'Enable user created by me'),
('user_enable_other', 'Enable user created by other'),
('user_update_itself', 'Update the user logged in only'),
('user_update_mine', 'Update user created by me'),
('user_update_other', 'Update user created by other');

CREATE TABLE IF NOT EXISTS publication (
id int(11) NOT NULL,
  title varchar(200) NOT NULL,
  year int(11) NOT NULL,
  book_title varchar(45) NOT NULL,
  chapter varchar(45) NOT NULL,
  city varchar(45) NOT NULL,
  country varchar(45) NOT NULL,
  conference_name varchar(45) NOT NULL,
  edition varchar(45) NOT NULL,
  institution varchar(45) NOT NULL,
  isbn varchar(45) NOT NULL,
  issn varchar(45) NOT NULL,
  journal varchar(45) NOT NULL,
  language varchar(45) NOT NULL,
  nationality varchar(45) NOT NULL,
  number varchar(45) NOT NULL,
  organization varchar(45) NOT NULL,
  pages varchar(45) NOT NULL,
  school varchar(45) NOT NULL COMMENT ' ',
  series varchar(45) NOT NULL,
  volume varchar(45) NOT NULL,
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL,
  publication_type int(11) NOT NULL,
  publisher int(11) NOT NULL,
  primary_author int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS publication_type (
id int(11) NOT NULL,
  name varchar(45) NOT NULL,
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS publisher (
id int(11) NOT NULL,
  name varchar(200) NOT NULL,
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS research_area (
id int(11) NOT NULL,
  name varchar(200) NOT NULL,
  logo varchar(32) NOT NULL DEFAULT '',
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS research_area_research_line (
  research_area int(11) NOT NULL,
  research_line int(11) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS research_line (
id int(11) NOT NULL,
  title varchar(200) NOT NULL,
  finished tinyint(1) NOT NULL,
  description text NOT NULL,
  logo varchar(32) NOT NULL DEFAULT '',
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL,
  primary_research_area int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS research_line_article (
  research_line int(11) NOT NULL,
  article int(11) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS research_line_financed_project (
  research_line int(11) NOT NULL,
  financed_project int(11) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS research_line_member (
  research_line int(11) NOT NULL,
  member int(11) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS research_line_partner (
  research_line int(11) NOT NULL,
  partner int(11) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS research_line_publication (
  research_line int(11) NOT NULL,
  publication int(11) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS research_line_resource (
  research_line int(11) NOT NULL,
  resource varchar(32) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS research_line_student_work (
  research_line int(11) NOT NULL,
  student_work int(11) NOT NULL,
  created_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS resource (
  filename_hash varchar(32) NOT NULL,
  filename varchar(45) NOT NULL,
  mime_type varchar(45) NOT NULL,
  size int(11) NOT NULL,
  private tinyint(1) NOT NULL,
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS rol (
  id varchar(45) NOT NULL,
  display_name varchar(45) NOT NULL,
  description varchar(200) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO rol (id, display_name, description) VALUES
('root_rol', 'Root rol', 'Root privileges');

CREATE TABLE IF NOT EXISTS rol_permission (
  rol varchar(45) NOT NULL,
  permission varchar(64) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS status (
id int(11) NOT NULL,
  name varchar(45) NOT NULL,
  description varchar(200) NOT NULL,
  created_by varchar(45) NOT NULL COMMENT ' ',
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS student_work (
id int(11) NOT NULL,
  title varchar(200) NOT NULL,
  year int(11) NOT NULL,
  school varchar(200) NOT NULL,
  volume varchar(45) NOT NULL,
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL,
  student_work_type int(11) NOT NULL,
  author int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS student_work_type (
id int(11) NOT NULL,
  name varchar(45) NOT NULL,
  created_by varchar(45) NOT NULL,
  updated_by varchar(45) NOT NULL,
  created_at int(11) NOT NULL,
  updated_at int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS timestamps (
  create_time timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ugroup (
  id varchar(45) NOT NULL,
  display_name varchar(45) NOT NULL,
  rol varchar(45) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO ugroup (id, display_name, rol) VALUES
('root_group', 'Root group', 'root_rol');

CREATE TABLE IF NOT EXISTS user (
  username varchar(45) NOT NULL,
  email varchar(45) NOT NULL,
  password varchar(60) NOT NULL,
  enabled tinyint(1) NOT NULL,
  display_name varchar(45) NOT NULL,
  ugroup varchar(45) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO user (username, email, password, enabled, display_name, ugroup) VALUES
('hugo', 'hugo@instanto.com', 'hugo', 1, 'Hugo González', 'root_group'),
('root', 'root@instanto.com', 'root', 1, 'Root user', 'root_group');


ALTER TABLE article
 ADD PRIMARY KEY (id), ADD KEY fk_article_1_idx (created_by), ADD KEY fk_article_2_idx (updated_by), ADD KEY fk_article_3_idx (newspaper);

ALTER TABLE category
 ADD PRIMARY KEY (id), ADD KEY fk_resource_type_1_idx (created_by), ADD KEY fk_resource_type_2_idx (updated_by);

ALTER TABLE category_resource
 ADD PRIMARY KEY (category,resource), ADD KEY fk_category_resource_1_idx (created_by), ADD KEY fk_category_resource_3_idx (resource);

ALTER TABLE config
 ADD PRIMARY KEY (config_key);

ALTER TABLE financed_project
 ADD PRIMARY KEY (id), ADD KEY fk_financed_project_1_idx (created_by), ADD KEY fk_financed_project_2_idx (updated_by), ADD KEY fk_financed_project_4_idx (primary_leader), ADD KEY fk_financed_project_3_idx (primary_funding_body);

ALTER TABLE financed_project_leader
 ADD PRIMARY KEY (financed_project,member), ADD KEY fk_financed_project_leader_1_idx (created_by), ADD KEY fk_financed_project_leader_3_idx (member);

ALTER TABLE financed_project_member
 ADD PRIMARY KEY (financed_project,member), ADD KEY fk_financed_project_member_1_idx (created_by), ADD KEY fk_financed_project_member_3_idx (member);

ALTER TABLE funding_body
 ADD PRIMARY KEY (id), ADD KEY fk_funding_body_1_idx (created_by), ADD KEY fk_funding_body_2_idx (updated_by);

ALTER TABLE funding_body_financed_project
 ADD PRIMARY KEY (funding_body,financed_project), ADD KEY fk_funding_body_financed_project_1_idx (created_by), ADD KEY fk_funding_body_financed_project_2_idx (updated_by), ADD KEY fk_funding_body_financed_project_4_idx (financed_project);

ALTER TABLE member
 ADD PRIMARY KEY (id), ADD KEY fk_member_1_idx (created_by), ADD KEY fk_member_2_idx (updated_by), ADD KEY fk_member_3_idx (primary_status);

ALTER TABLE member_publication
 ADD PRIMARY KEY (member,publication), ADD KEY fk_member_publication_1_idx (created_by), ADD KEY fk_member_publication_3_idx (publication);

ALTER TABLE member_status
 ADD PRIMARY KEY (member,status), ADD KEY fk_member_status_1_idx (created_by), ADD KEY fk_member_status_3_idx (status);

ALTER TABLE newspaper
 ADD PRIMARY KEY (id), ADD KEY fk_newspaper_1_idx (created_by), ADD KEY fk_newspaper_2_idx (updated_by);

ALTER TABLE partner
 ADD PRIMARY KEY (id), ADD KEY fk_partner_1_idx (created_by), ADD KEY fk_partner_2_idx (updated_by);

ALTER TABLE partner_member
 ADD PRIMARY KEY (partner,member), ADD KEY fk_partner_member_1_idx (created_by), ADD KEY fk_partner_member_3_idx (member);

ALTER TABLE permission
 ADD PRIMARY KEY (id);

ALTER TABLE publication
 ADD PRIMARY KEY (id), ADD KEY fk_publication_1_idx (created_by), ADD KEY fk_publication_2_idx (updated_by), ADD KEY fk_publication_3_idx (publication_type), ADD KEY fk_publication_4_idx (publisher), ADD KEY fk_publication_5_idx (primary_author);

ALTER TABLE publication_type
 ADD PRIMARY KEY (id), ADD KEY fk_publication_type_1_idx (created_by), ADD KEY fk_publication_type_2_idx (updated_by);

ALTER TABLE publisher
 ADD PRIMARY KEY (id), ADD KEY fk_publisher_1_idx (created_by), ADD KEY fk_publisher_2_idx (updated_by);

ALTER TABLE research_area
 ADD PRIMARY KEY (id), ADD KEY fk_research_area_1_idx (created_by), ADD KEY fk_research_area_2_idx (updated_by);

ALTER TABLE research_area_research_line
 ADD PRIMARY KEY (research_area,research_line), ADD KEY fk_research_area_research_line_1_idx (created_by), ADD KEY fk_research_area_research_line_3_idx (research_line);

ALTER TABLE research_line
 ADD PRIMARY KEY (id), ADD KEY fk_research_line_1_idx (created_by), ADD KEY fk_research_line_2_idx (updated_by), ADD KEY fk_research_line_3_idx (primary_research_area);

ALTER TABLE research_line_article
 ADD PRIMARY KEY (research_line,article), ADD KEY fk_research_line_article_1_idx (created_by), ADD KEY fk_research_line_article_3_idx (article);

ALTER TABLE research_line_financed_project
 ADD PRIMARY KEY (research_line,financed_project), ADD KEY fk_research_line_financed_project_1_idx (created_by), ADD KEY fk_research_line_financed_project_3_idx (financed_project);

ALTER TABLE research_line_member
 ADD PRIMARY KEY (research_line,member), ADD KEY fk_research_line_member_1_idx (created_by), ADD KEY fk_research_line_member_3_idx (member);

ALTER TABLE research_line_partner
 ADD PRIMARY KEY (research_line,partner), ADD KEY fk_research_line_partner_1_idx (created_by), ADD KEY fk_research_line_partner_3_idx (partner);

ALTER TABLE research_line_publication
 ADD PRIMARY KEY (research_line,publication), ADD KEY fk_research_line_publication_1_idx (created_by), ADD KEY fk_research_line_publication_3_idx (publication);

ALTER TABLE research_line_resource
 ADD PRIMARY KEY (research_line,resource), ADD KEY fk_research_line_resource_1_idx (created_by), ADD KEY fk_research_line_resource_3_idx (resource);

ALTER TABLE research_line_student_work
 ADD PRIMARY KEY (research_line,student_work), ADD KEY fk_research_line_student_work_1_idx (created_by), ADD KEY fk_research_line_student_work_3_idx (student_work);

ALTER TABLE resource
 ADD PRIMARY KEY (filename_hash), ADD KEY fk_resource_1_idx (created_by), ADD KEY fk_resource_2_idx (updated_by);

ALTER TABLE rol
 ADD PRIMARY KEY (id);

ALTER TABLE rol_permission
 ADD PRIMARY KEY (rol,permission), ADD KEY fk_rol_permission_3_idx (permission);

ALTER TABLE status
 ADD PRIMARY KEY (id), ADD KEY fk_status_1_idx (created_by), ADD KEY fk_status_2_idx (updated_by);

ALTER TABLE student_work
 ADD PRIMARY KEY (id), ADD KEY fk_student_work_1_idx (created_by), ADD KEY fk_student_work_2_idx (updated_by), ADD KEY fk_student_work_3_idx (student_work_type), ADD KEY fk_student_work_4_idx (author);

ALTER TABLE student_work_type
 ADD PRIMARY KEY (id), ADD KEY fk_student_work_type_1_idx (created_by), ADD KEY fk_student_work_type_2_idx (updated_by);

ALTER TABLE ugroup
 ADD PRIMARY KEY (id), ADD KEY fk_ugroup_rol_idx (rol);

ALTER TABLE user
 ADD PRIMARY KEY (username), ADD KEY fk_user_ugroup_ugroup_idx (ugroup);


ALTER TABLE article
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE category
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE financed_project
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE funding_body
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE member
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE newspaper
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE partner
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE publication
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE publication_type
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE publisher
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE research_area
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE research_line
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE status
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE student_work
MODIFY id int(11) NOT NULL AUTO_INCREMENT;
ALTER TABLE student_work_type
MODIFY id int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE article
ADD CONSTRAINT fk_article_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_article_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_article_3 FOREIGN KEY (newspaper) REFERENCES newspaper (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE category
ADD CONSTRAINT fk_resource_type_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_resource_type_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE category_resource
ADD CONSTRAINT fk_category_resource_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_category_resource_2 FOREIGN KEY (category) REFERENCES category (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_category_resource_3 FOREIGN KEY (resource) REFERENCES resource (filename_hash) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE financed_project
ADD CONSTRAINT fk_financed_project_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_financed_project_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_financed_project_3 FOREIGN KEY (primary_funding_body) REFERENCES funding_body (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_financed_project_4 FOREIGN KEY (primary_leader) REFERENCES member (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE financed_project_leader
ADD CONSTRAINT fk_financed_project_leader_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_financed_project_leader_2 FOREIGN KEY (financed_project) REFERENCES financed_project (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_financed_project_leader_3 FOREIGN KEY (member) REFERENCES member (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE financed_project_member
ADD CONSTRAINT fk_financed_project_member_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_financed_project_member_2 FOREIGN KEY (financed_project) REFERENCES financed_project (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_financed_project_member_3 FOREIGN KEY (member) REFERENCES member (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE funding_body
ADD CONSTRAINT fk_funding_body_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_funding_body_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE funding_body_financed_project
ADD CONSTRAINT fk_funding_body_financed_project_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_funding_body_financed_project_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_funding_body_financed_project_3 FOREIGN KEY (funding_body) REFERENCES funding_body (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_funding_body_financed_project_4 FOREIGN KEY (financed_project) REFERENCES financed_project (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE member
ADD CONSTRAINT fk_member_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_member_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_member_3 FOREIGN KEY (primary_status) REFERENCES status (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE member_publication
ADD CONSTRAINT fk_member_publication_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_member_publication_2 FOREIGN KEY (member) REFERENCES member (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_member_publication_3 FOREIGN KEY (publication) REFERENCES publication (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE member_status
ADD CONSTRAINT fk_member_status_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_member_status_2 FOREIGN KEY (member) REFERENCES member (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_member_status_3 FOREIGN KEY (status) REFERENCES status (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE newspaper
ADD CONSTRAINT fk_newspaper_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_newspaper_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE partner
ADD CONSTRAINT fk_partner_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_partner_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE partner_member
ADD CONSTRAINT fk_partner_member_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_partner_member_2 FOREIGN KEY (partner) REFERENCES partner (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_partner_member_3 FOREIGN KEY (member) REFERENCES member (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE publication
ADD CONSTRAINT fk_publication_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_publication_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_publication_3 FOREIGN KEY (publication_type) REFERENCES publication_type (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_publication_4 FOREIGN KEY (publisher) REFERENCES publisher (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_publication_5 FOREIGN KEY (primary_author) REFERENCES member (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE publication_type
ADD CONSTRAINT fk_publication_type_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_publication_type_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE publisher
ADD CONSTRAINT fk_publisher_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_publisher_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE research_area
ADD CONSTRAINT fk_research_area_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_research_area_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE research_area_research_line
ADD CONSTRAINT fk_research_area_research_line_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_research_area_research_line_2 FOREIGN KEY (research_area) REFERENCES research_area (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_research_area_research_line_3 FOREIGN KEY (research_line) REFERENCES research_line (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE research_line
ADD CONSTRAINT fk_research_line_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_research_line_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_research_line_3 FOREIGN KEY (primary_research_area) REFERENCES research_area (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE research_line_article
ADD CONSTRAINT fk_research_line_article_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_research_line_article_2 FOREIGN KEY (research_line) REFERENCES research_line (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_research_line_article_3 FOREIGN KEY (article) REFERENCES article (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE research_line_financed_project
ADD CONSTRAINT fk_research_line_financed_project_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_research_line_financed_project_2 FOREIGN KEY (research_line) REFERENCES research_line (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_research_line_financed_project_3 FOREIGN KEY (financed_project) REFERENCES financed_project (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE research_line_member
ADD CONSTRAINT fk_research_line_member_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_research_line_member_2 FOREIGN KEY (research_line) REFERENCES research_line (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_research_line_member_3 FOREIGN KEY (member) REFERENCES member (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE research_line_partner
ADD CONSTRAINT fk_research_line_partner_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_research_line_partner_2 FOREIGN KEY (research_line) REFERENCES research_line (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_research_line_partner_3 FOREIGN KEY (partner) REFERENCES partner (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE research_line_publication
ADD CONSTRAINT fk_research_line_publication_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_research_line_publication_2 FOREIGN KEY (research_line) REFERENCES research_line (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_research_line_publication_3 FOREIGN KEY (publication) REFERENCES publication (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE research_line_resource
ADD CONSTRAINT fk_research_line_resource_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_research_line_resource_2 FOREIGN KEY (research_line) REFERENCES research_line (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_research_line_resource_3 FOREIGN KEY (resource) REFERENCES resource (filename_hash) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE research_line_student_work
ADD CONSTRAINT fk_research_line_student_work_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_research_line_student_work_2 FOREIGN KEY (research_line) REFERENCES research_line (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_research_line_student_work_3 FOREIGN KEY (student_work) REFERENCES student_work (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE resource
ADD CONSTRAINT fk_resource_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_resource_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE rol_permission
ADD CONSTRAINT fk_rol_permission_2 FOREIGN KEY (rol) REFERENCES rol (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_rol_permission_3 FOREIGN KEY (permission) REFERENCES permission (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE status
ADD CONSTRAINT fk_status_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_status_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE student_work
ADD CONSTRAINT fk_student_work_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_student_work_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_student_work_3 FOREIGN KEY (student_work_type) REFERENCES student_work_type (id) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT fk_student_work_4 FOREIGN KEY (author) REFERENCES member (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE student_work_type
ADD CONSTRAINT fk_student_work_type_1 FOREIGN KEY (created_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION,
ADD CONSTRAINT fk_student_work_type_2 FOREIGN KEY (updated_by) REFERENCES user (username) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE ugroup
ADD CONSTRAINT fk_ugroup_rol FOREIGN KEY (rol) REFERENCES rol (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE user
ADD CONSTRAINT fk_user_ugroup_ugroup FOREIGN KEY (ugroup) REFERENCES ugroup (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
`

	db, err := sqlx.Connect("mysql", "uranus:uranus@tcp(192.168.31.135:3306)/instanto")
	if err != nil {
		fmt.Println(err)
		return
	}

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	fmt.Println(schema)
	db.MustExec("CREATE DATABASE hello")

}
