--
-- PostgreSQL database dump
--

\restrict RTsRvWeYoU4RHyP6dxElvXPuDz6ji7fjTT5Rb0koqeL8bIYvPWuZimWhqp5a0Bb

-- Dumped from database version 16.9
-- Dumped by pg_dump version 18.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE IF EXISTS ONLY public.tournee DROP CONSTRAINT IF EXISTS tournee_id_planning_fkey;
ALTER TABLE IF EXISTS ONLY public.tournee DROP CONSTRAINT IF EXISTS tournee_id_destinataire_fkey;
ALTER TABLE IF EXISTS ONLY public.tournee DROP CONSTRAINT IF EXISTS tournee_id_benevole_fkey;
ALTER TABLE IF EXISTS ONLY public.tournee DROP CONSTRAINT IF EXISTS tournee_id_agence_fkey;
ALTER TABLE IF EXISTS ONLY public.stock DROP CONSTRAINT IF EXISTS stock_id_agence_fkey;
ALTER TABLE IF EXISTS ONLY public.service DROP CONSTRAINT IF EXISTS service_id_competence_fkey;
ALTER TABLE IF EXISTS ONLY public.service DROP CONSTRAINT IF EXISTS service_id_agence_fkey;
ALTER TABLE IF EXISTS ONLY public.service_competence DROP CONSTRAINT IF EXISTS service_competence_id_service_fkey;
ALTER TABLE IF EXISTS ONLY public.service_competence DROP CONSTRAINT IF EXISTS service_competence_id_competence_fkey;
ALTER TABLE IF EXISTS ONLY public.produit_tournee DROP CONSTRAINT IF EXISTS produit_tournee_id_tournee_fkey;
ALTER TABLE IF EXISTS ONLY public.produit_tournee DROP CONSTRAINT IF EXISTS produit_tournee_id_stock_fkey;
ALTER TABLE IF EXISTS ONLY public.produit_collecte DROP CONSTRAINT IF EXISTS produit_collecte_id_stock_fkey;
ALTER TABLE IF EXISTS ONLY public.produit_collecte DROP CONSTRAINT IF EXISTS produit_collecte_id_collecte_fkey;
ALTER TABLE IF EXISTS ONLY public.planning_excel DROP CONSTRAINT IF EXISTS planning_excel_id_planning_fkey;
ALTER TABLE IF EXISTS ONLY public.planning_excel DROP CONSTRAINT IF EXISTS planning_excel_id_benevole_fkey;
ALTER TABLE IF EXISTS ONLY public.commercant DROP CONSTRAINT IF EXISTS fk_utilisateur_commercant;
ALTER TABLE IF EXISTS ONLY public.benevole DROP CONSTRAINT IF EXISTS fk_utilisateur_benevole;
ALTER TABLE IF EXISTS ONLY public.association_beneficiaire DROP CONSTRAINT IF EXISTS fk_utilisateur_association;
ALTER TABLE IF EXISTS ONLY public.adherent DROP CONSTRAINT IF EXISTS fk_utilisateur_adherent;
ALTER TABLE IF EXISTS ONLY public.utilisateur DROP CONSTRAINT IF EXISTS fk_role;
ALTER TABLE IF EXISTS ONLY public.disponibilite DROP CONSTRAINT IF EXISTS fk_disponibilite_benevole;
ALTER TABLE IF EXISTS ONLY public.planning DROP CONSTRAINT IF EXISTS fk_disponibilite;
ALTER TABLE IF EXISTS ONLY public.demande_service DROP CONSTRAINT IF EXISTS fk_demande_service;
ALTER TABLE IF EXISTS ONLY public.demande_service DROP CONSTRAINT IF EXISTS fk_demande_adherent;
ALTER TABLE IF EXISTS ONLY public.benevole_competence DROP CONSTRAINT IF EXISTS fk_competence;
ALTER TABLE IF EXISTS ONLY public.planning DROP CONSTRAINT IF EXISTS fk_benevole;
ALTER TABLE IF EXISTS ONLY public.benevole_competence DROP CONSTRAINT IF EXISTS fk_benevole;
ALTER TABLE IF EXISTS ONLY public.utilisateur DROP CONSTRAINT IF EXISTS fk_agence;
ALTER TABLE IF EXISTS ONLY public.cotisation DROP CONSTRAINT IF EXISTS fk_adherent_cotisation;
ALTER TABLE IF EXISTS ONLY public.destinataire DROP CONSTRAINT IF EXISTS destinataire_id_agence_fkey;
ALTER TABLE IF EXISTS ONLY public.demande_service DROP CONSTRAINT IF EXISTS demande_service_id_planning_fkey;
ALTER TABLE IF EXISTS ONLY public.demande_service DROP CONSTRAINT IF EXISTS demande_service_id_benevole_fkey;
ALTER TABLE IF EXISTS ONLY public.collecte DROP CONSTRAINT IF EXISTS collecte_id_planning_fkey;
ALTER TABLE IF EXISTS ONLY public.collecte DROP CONSTRAINT IF EXISTS collecte_id_commercant_fkey;
ALTER TABLE IF EXISTS ONLY public.collecte DROP CONSTRAINT IF EXISTS collecte_id_benevole_fkey;
ALTER TABLE IF EXISTS ONLY public.collecte DROP CONSTRAINT IF EXISTS collecte_id_agence_fkey;
ALTER TABLE IF EXISTS ONLY public.benevole DROP CONSTRAINT IF EXISTS benevole_id_competence_fkey;
ALTER TABLE IF EXISTS ONLY public.utilisateur DROP CONSTRAINT IF EXISTS utilisateur_pkey;
ALTER TABLE IF EXISTS ONLY public.utilisateur DROP CONSTRAINT IF EXISTS utilisateur_email_key;
ALTER TABLE IF EXISTS ONLY public.tournee DROP CONSTRAINT IF EXISTS tournee_pkey;
ALTER TABLE IF EXISTS ONLY public.stock DROP CONSTRAINT IF EXISTS stock_pkey;
ALTER TABLE IF EXISTS ONLY public.service DROP CONSTRAINT IF EXISTS service_pkey;
ALTER TABLE IF EXISTS ONLY public.service_competence DROP CONSTRAINT IF EXISTS service_competence_pkey;
ALTER TABLE IF EXISTS ONLY public.role DROP CONSTRAINT IF EXISTS role_pkey;
ALTER TABLE IF EXISTS ONLY public.role DROP CONSTRAINT IF EXISTS role_nom_key;
ALTER TABLE IF EXISTS ONLY public.produit_tournee DROP CONSTRAINT IF EXISTS produit_tournee_pkey;
ALTER TABLE IF EXISTS ONLY public.produit_collecte DROP CONSTRAINT IF EXISTS produit_collecte_pkey;
ALTER TABLE IF EXISTS ONLY public.produit_collecte DROP CONSTRAINT IF EXISTS produit_collecte_code_barre_key;
ALTER TABLE IF EXISTS ONLY public.planning DROP CONSTRAINT IF EXISTS planning_pkey;
ALTER TABLE IF EXISTS ONLY public.planning_excel DROP CONSTRAINT IF EXISTS planning_excel_pkey;
ALTER TABLE IF EXISTS ONLY public.disponibilite DROP CONSTRAINT IF EXISTS disponibilite_pkey;
ALTER TABLE IF EXISTS ONLY public.destinataire DROP CONSTRAINT IF EXISTS destinataire_pkey;
ALTER TABLE IF EXISTS ONLY public.demande_service DROP CONSTRAINT IF EXISTS demande_service_pkey;
ALTER TABLE IF EXISTS ONLY public.cotisation DROP CONSTRAINT IF EXISTS cotisation_stripe_session_id_key;
ALTER TABLE IF EXISTS ONLY public.cotisation DROP CONSTRAINT IF EXISTS cotisation_pkey;
ALTER TABLE IF EXISTS ONLY public.competence DROP CONSTRAINT IF EXISTS competence_pkey;
ALTER TABLE IF EXISTS ONLY public.commercant DROP CONSTRAINT IF EXISTS commercant_pkey;
ALTER TABLE IF EXISTS ONLY public.commercant DROP CONSTRAINT IF EXISTS commercant_id_utilisateur_key;
ALTER TABLE IF EXISTS ONLY public.collecte DROP CONSTRAINT IF EXISTS collecte_pkey;
ALTER TABLE IF EXISTS ONLY public.benevole DROP CONSTRAINT IF EXISTS benevole_pkey;
ALTER TABLE IF EXISTS ONLY public.benevole DROP CONSTRAINT IF EXISTS benevole_id_utilisateur_key;
ALTER TABLE IF EXISTS ONLY public.benevole_competence DROP CONSTRAINT IF EXISTS benevole_competence_pkey;
ALTER TABLE IF EXISTS ONLY public.association_beneficiaire DROP CONSTRAINT IF EXISTS association_beneficiaire_pkey;
ALTER TABLE IF EXISTS ONLY public.association_beneficiaire DROP CONSTRAINT IF EXISTS association_beneficiaire_id_utilisateur_key;
ALTER TABLE IF EXISTS ONLY public.agence DROP CONSTRAINT IF EXISTS agence_pkey;
ALTER TABLE IF EXISTS ONLY public.adherent DROP CONSTRAINT IF EXISTS adherent_pkey;
ALTER TABLE IF EXISTS ONLY public.adherent DROP CONSTRAINT IF EXISTS adherent_id_utilisateur_key;
ALTER TABLE IF EXISTS public.utilisateur ALTER COLUMN id_utilisateur DROP DEFAULT;
ALTER TABLE IF EXISTS public.tournee ALTER COLUMN id_tournee DROP DEFAULT;
ALTER TABLE IF EXISTS public.stock ALTER COLUMN id_stock DROP DEFAULT;
ALTER TABLE IF EXISTS public.service ALTER COLUMN id_service DROP DEFAULT;
ALTER TABLE IF EXISTS public.role ALTER COLUMN id_role DROP DEFAULT;
ALTER TABLE IF EXISTS public.produit_tournee ALTER COLUMN id_produit_tournee DROP DEFAULT;
ALTER TABLE IF EXISTS public.produit_collecte ALTER COLUMN id_produit_collecte DROP DEFAULT;
ALTER TABLE IF EXISTS public.planning_excel ALTER COLUMN id_planning_excel DROP DEFAULT;
ALTER TABLE IF EXISTS public.planning ALTER COLUMN id_planning DROP DEFAULT;
ALTER TABLE IF EXISTS public.disponibilite ALTER COLUMN id_disponibilite DROP DEFAULT;
ALTER TABLE IF EXISTS public.destinataire ALTER COLUMN id_destinataire DROP DEFAULT;
ALTER TABLE IF EXISTS public.demande_service ALTER COLUMN id_demande_service DROP DEFAULT;
ALTER TABLE IF EXISTS public.cotisation ALTER COLUMN id_cotisation DROP DEFAULT;
ALTER TABLE IF EXISTS public.competence ALTER COLUMN id_competence DROP DEFAULT;
ALTER TABLE IF EXISTS public.commercant ALTER COLUMN id_commercant DROP DEFAULT;
ALTER TABLE IF EXISTS public.collecte ALTER COLUMN id_collecte DROP DEFAULT;
ALTER TABLE IF EXISTS public.benevole ALTER COLUMN id_benevole DROP DEFAULT;
ALTER TABLE IF EXISTS public.association_beneficiaire ALTER COLUMN id_association DROP DEFAULT;
ALTER TABLE IF EXISTS public.agence ALTER COLUMN id_agence DROP DEFAULT;
ALTER TABLE IF EXISTS public.adherent ALTER COLUMN id_adherent DROP DEFAULT;
DROP SEQUENCE IF EXISTS public.utilisateur_id_utilisateur_seq;
DROP TABLE IF EXISTS public.utilisateur;
DROP SEQUENCE IF EXISTS public.tournee_id_tournee_seq;
DROP TABLE IF EXISTS public.tournee;
DROP SEQUENCE IF EXISTS public.stock_id_stock_seq;
DROP TABLE IF EXISTS public.stock;
DROP SEQUENCE IF EXISTS public.service_id_service_seq;
DROP TABLE IF EXISTS public.service_competence;
DROP TABLE IF EXISTS public.service;
DROP SEQUENCE IF EXISTS public.role_id_role_seq;
DROP TABLE IF EXISTS public.role;
DROP SEQUENCE IF EXISTS public.produit_tournee_id_produit_tournee_seq;
DROP TABLE IF EXISTS public.produit_tournee;
DROP SEQUENCE IF EXISTS public.produit_collecte_id_produit_collecte_seq;
DROP TABLE IF EXISTS public.produit_collecte;
DROP SEQUENCE IF EXISTS public.planning_id_planning_seq;
DROP SEQUENCE IF EXISTS public.planning_excel_id_planning_excel_seq;
DROP TABLE IF EXISTS public.planning_excel;
DROP TABLE IF EXISTS public.planning;
DROP SEQUENCE IF EXISTS public.disponibilite_id_disponibilite_seq;
DROP TABLE IF EXISTS public.disponibilite;
DROP SEQUENCE IF EXISTS public.destinataire_id_destinataire_seq;
DROP TABLE IF EXISTS public.destinataire;
DROP SEQUENCE IF EXISTS public.demande_service_id_demande_service_seq;
DROP TABLE IF EXISTS public.demande_service;
DROP SEQUENCE IF EXISTS public.cotisation_id_cotisation_seq;
DROP TABLE IF EXISTS public.cotisation;
DROP SEQUENCE IF EXISTS public.competence_id_competence_seq;
DROP TABLE IF EXISTS public.competence;
DROP SEQUENCE IF EXISTS public.commercant_id_commercant_seq;
DROP TABLE IF EXISTS public.commercant;
DROP SEQUENCE IF EXISTS public.collecte_id_collecte_seq;
DROP TABLE IF EXISTS public.collecte;
DROP SEQUENCE IF EXISTS public.benevole_id_benevole_seq;
DROP TABLE IF EXISTS public.benevole_competence;
DROP TABLE IF EXISTS public.benevole;
DROP SEQUENCE IF EXISTS public.association_beneficiaire_id_association_seq;
DROP TABLE IF EXISTS public.association_beneficiaire;
DROP SEQUENCE IF EXISTS public.agence_id_agence_seq;
DROP TABLE IF EXISTS public.agence;
DROP SEQUENCE IF EXISTS public.adherent_id_adherent_seq;
DROP TABLE IF EXISTS public.adherent;
DROP TYPE IF EXISTS public.type_destinataire;
DROP TYPE IF EXISTS public.statut_tournee;
DROP TYPE IF EXISTS public.statut_collecte;
--
-- Name: statut_collecte; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.statut_collecte AS ENUM (
    'en_attente',
    'planifiee',
    'effectuee'
);


ALTER TYPE public.statut_collecte OWNER TO postgres;

--
-- Name: statut_tournee; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.statut_tournee AS ENUM (
    'en_attente',
    'planifiee',
    'effectuee'
);


ALTER TYPE public.statut_tournee OWNER TO postgres;

--
-- Name: type_destinataire; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.type_destinataire AS ENUM (
    'ASSOCIATION',
    'PARTICULIER'
);


ALTER TYPE public.type_destinataire OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: adherent; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.adherent (
    id_adherent integer NOT NULL,
    id_utilisateur integer NOT NULL,
    date_adhesion date,
    date_expiration date,
    statut character varying(20) NOT NULL,
    CONSTRAINT adherent_statut_check CHECK (((statut)::text = ANY ((ARRAY['EN_ATTENTE'::character varying, 'ACTIF'::character varying, 'EXPIRE'::character varying, 'SUSPENDU'::character varying])::text[])))
);


ALTER TABLE public.adherent OWNER TO postgres;

--
-- Name: adherent_id_adherent_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.adherent_id_adherent_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.adherent_id_adherent_seq OWNER TO postgres;

--
-- Name: adherent_id_adherent_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.adherent_id_adherent_seq OWNED BY public.adherent.id_adherent;


--
-- Name: agence; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.agence (
    id_agence integer NOT NULL,
    nom character varying(100) NOT NULL,
    adresse character varying(255),
    ville character varying(100),
    code_postal character varying(10),
    pays character varying(100),
    telephone character varying(20),
    email character varying(150)
);


ALTER TABLE public.agence OWNER TO postgres;

--
-- Name: agence_id_agence_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.agence_id_agence_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.agence_id_agence_seq OWNER TO postgres;

--
-- Name: agence_id_agence_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.agence_id_agence_seq OWNED BY public.agence.id_agence;


--
-- Name: association_beneficiaire; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.association_beneficiaire (
    id_association integer NOT NULL,
    id_utilisateur integer NOT NULL,
    nom_responsable character varying(150),
    nom_association character varying(150),
    nombre_beneficiaires integer,
    type_association character varying(100),
    date_validation date,
    statut character varying(20) DEFAULT 'EN_ATTENTE'::character varying,
    CONSTRAINT association_beneficiaire_statut_check CHECK (((statut)::text = ANY ((ARRAY['EN_ATTENTE'::character varying, 'VALIDEE'::character varying, 'REFUSEE'::character varying])::text[])))
);


ALTER TABLE public.association_beneficiaire OWNER TO postgres;

--
-- Name: association_beneficiaire_id_association_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.association_beneficiaire_id_association_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.association_beneficiaire_id_association_seq OWNER TO postgres;

--
-- Name: association_beneficiaire_id_association_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.association_beneficiaire_id_association_seq OWNED BY public.association_beneficiaire.id_association;


--
-- Name: benevole; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.benevole (
    id_benevole integer NOT NULL,
    id_utilisateur integer NOT NULL,
    permis boolean DEFAULT false,
    disponibilite text,
    statut character varying(20) DEFAULT 'EN_ATTENTE'::character varying,
    id_competence integer,
    justificatif character varying(255),
    commentaire_validation text,
    CONSTRAINT benevole_statut_check CHECK (((statut)::text = ANY ((ARRAY['EN_ATTENTE'::character varying, 'VALIDE'::character varying, 'REFUSE'::character varying])::text[])))
);


ALTER TABLE public.benevole OWNER TO postgres;

--
-- Name: benevole_competence; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.benevole_competence (
    id_benevole integer NOT NULL,
    id_competence integer NOT NULL
);


ALTER TABLE public.benevole_competence OWNER TO postgres;

--
-- Name: benevole_id_benevole_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.benevole_id_benevole_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.benevole_id_benevole_seq OWNER TO postgres;

--
-- Name: benevole_id_benevole_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.benevole_id_benevole_seq OWNED BY public.benevole.id_benevole;


--
-- Name: collecte; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.collecte (
    id_collecte integer NOT NULL,
    id_commercant integer NOT NULL,
    id_agence integer NOT NULL,
    id_benevole integer,
    id_planning integer,
    date_collecte date NOT NULL,
    statut public.statut_collecte DEFAULT 'en_attente'::public.statut_collecte NOT NULL,
    pdf_recapitulatif character varying(255)
);


ALTER TABLE public.collecte OWNER TO postgres;

--
-- Name: collecte_id_collecte_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.collecte_id_collecte_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.collecte_id_collecte_seq OWNER TO postgres;

--
-- Name: collecte_id_collecte_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.collecte_id_collecte_seq OWNED BY public.collecte.id_collecte;


--
-- Name: commercant; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.commercant (
    id_commercant integer NOT NULL,
    id_utilisateur integer NOT NULL,
    nom_entreprise character varying(150) NOT NULL,
    type_commerce character varying(100),
    numero_siret character varying(20),
    date_adhesion date,
    date_expiration date,
    cotisation numeric(10,2),
    statut character varying(20) DEFAULT 'ACTIF'::character varying,
    CONSTRAINT commercant_statut_check CHECK (((statut)::text = ANY ((ARRAY['ACTIF'::character varying, 'EXPIRE'::character varying, 'SUSPENDU'::character varying])::text[])))
);


ALTER TABLE public.commercant OWNER TO postgres;

--
-- Name: commercant_id_commercant_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.commercant_id_commercant_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.commercant_id_commercant_seq OWNER TO postgres;

--
-- Name: commercant_id_commercant_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.commercant_id_commercant_seq OWNED BY public.commercant.id_commercant;


--
-- Name: competence; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.competence (
    id_competence integer NOT NULL,
    nom character varying(100) NOT NULL,
    description text
);


ALTER TABLE public.competence OWNER TO postgres;

--
-- Name: competence_id_competence_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.competence_id_competence_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.competence_id_competence_seq OWNER TO postgres;

--
-- Name: competence_id_competence_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.competence_id_competence_seq OWNED BY public.competence.id_competence;


--
-- Name: cotisation; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cotisation (
    id_cotisation integer NOT NULL,
    id_adherent integer NOT NULL,
    montant numeric(10,2) NOT NULL,
    date_paiement date,
    date_debut date,
    date_expiration date,
    statut character varying(20) NOT NULL,
    stripe_session_id character varying(255),
    CONSTRAINT cotisation_statut_check CHECK (((statut)::text = ANY ((ARRAY['EN_ATTENTE'::character varying, 'PAYEE'::character varying, 'ANNULEE'::character varying])::text[])))
);


ALTER TABLE public.cotisation OWNER TO postgres;

--
-- Name: cotisation_id_cotisation_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cotisation_id_cotisation_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.cotisation_id_cotisation_seq OWNER TO postgres;

--
-- Name: cotisation_id_cotisation_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cotisation_id_cotisation_seq OWNED BY public.cotisation.id_cotisation;


--
-- Name: demande_service; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.demande_service (
    id_demande_service integer NOT NULL,
    id_service integer NOT NULL,
    id_adherent integer NOT NULL,
    date_demande date DEFAULT CURRENT_DATE NOT NULL,
    statut character varying(20) DEFAULT 'EN_ATTENTE'::character varying NOT NULL,
    date_creation timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    id_benevole integer,
    id_planning integer,
    CONSTRAINT chk_statut_demande CHECK (((statut)::text = ANY ((ARRAY['EN_ATTENTE'::character varying, 'ATTRIBUE'::character varying, 'TERMINE'::character varying, 'ANNULE'::character varying])::text[])))
);


ALTER TABLE public.demande_service OWNER TO postgres;

--
-- Name: demande_service_id_demande_service_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.demande_service_id_demande_service_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.demande_service_id_demande_service_seq OWNER TO postgres;

--
-- Name: demande_service_id_demande_service_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.demande_service_id_demande_service_seq OWNED BY public.demande_service.id_demande_service;


--
-- Name: destinataire; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.destinataire (
    id_destinataire integer NOT NULL,
    id_agence integer NOT NULL,
    type public.type_destinataire NOT NULL,
    nom character varying(150) NOT NULL,
    adresse character varying(255) NOT NULL
);


ALTER TABLE public.destinataire OWNER TO postgres;

--
-- Name: destinataire_id_destinataire_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.destinataire_id_destinataire_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.destinataire_id_destinataire_seq OWNER TO postgres;

--
-- Name: destinataire_id_destinataire_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.destinataire_id_destinataire_seq OWNED BY public.destinataire.id_destinataire;


--
-- Name: disponibilite; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.disponibilite (
    id_disponibilite integer NOT NULL,
    id_benevole integer NOT NULL,
    date_disponibilite date NOT NULL,
    heure_debut time without time zone NOT NULL,
    heure_fin time without time zone NOT NULL,
    statut character varying(20) DEFAULT 'DISPONIBLE'::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_heure CHECK ((heure_fin > heure_debut)),
    CONSTRAINT chk_statut CHECK (((statut)::text = ANY ((ARRAY['DISPONIBLE'::character varying, 'ABSENT'::character varying, 'RESERVE'::character varying])::text[])))
);


ALTER TABLE public.disponibilite OWNER TO postgres;

--
-- Name: disponibilite_id_disponibilite_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.disponibilite_id_disponibilite_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.disponibilite_id_disponibilite_seq OWNER TO postgres;

--
-- Name: disponibilite_id_disponibilite_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.disponibilite_id_disponibilite_seq OWNED BY public.disponibilite.id_disponibilite;


--
-- Name: planning; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.planning (
    id_planning integer NOT NULL,
    id_benevole integer NOT NULL,
    id_disponibilite integer NOT NULL,
    date date NOT NULL,
    heure_debut time without time zone NOT NULL,
    heure_fin time without time zone NOT NULL,
    statut character varying(20) DEFAULT 'PLANIFIE'::character varying,
    CONSTRAINT planning_statut_check CHECK (((statut)::text = ANY ((ARRAY['PLANIFIE'::character varying, 'ATTRIBUE'::character varying, 'TERMINE'::character varying, 'ANNULE'::character varying])::text[])))
);


ALTER TABLE public.planning OWNER TO postgres;

--
-- Name: planning_excel; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.planning_excel (
    id_planning_excel integer NOT NULL,
    id_planning integer NOT NULL,
    id_benevole integer NOT NULL,
    chemin_fichier text NOT NULL,
    date_generation timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.planning_excel OWNER TO postgres;

--
-- Name: planning_excel_id_planning_excel_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.planning_excel_id_planning_excel_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.planning_excel_id_planning_excel_seq OWNER TO postgres;

--
-- Name: planning_excel_id_planning_excel_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.planning_excel_id_planning_excel_seq OWNED BY public.planning_excel.id_planning_excel;


--
-- Name: planning_id_planning_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.planning_id_planning_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.planning_id_planning_seq OWNER TO postgres;

--
-- Name: planning_id_planning_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.planning_id_planning_seq OWNED BY public.planning.id_planning;


--
-- Name: produit_collecte; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.produit_collecte (
    id_produit_collecte integer NOT NULL,
    id_collecte integer NOT NULL,
    libelle character varying(150) NOT NULL,
    quantite numeric(10,2) NOT NULL,
    code_barre character varying(50),
    id_stock integer
);


ALTER TABLE public.produit_collecte OWNER TO postgres;

--
-- Name: produit_collecte_id_produit_collecte_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.produit_collecte_id_produit_collecte_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.produit_collecte_id_produit_collecte_seq OWNER TO postgres;

--
-- Name: produit_collecte_id_produit_collecte_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.produit_collecte_id_produit_collecte_seq OWNED BY public.produit_collecte.id_produit_collecte;


--
-- Name: produit_tournee; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.produit_tournee (
    id_produit_tournee integer NOT NULL,
    id_tournee integer NOT NULL,
    id_stock integer NOT NULL,
    quantite numeric(10,2) NOT NULL
);


ALTER TABLE public.produit_tournee OWNER TO postgres;

--
-- Name: produit_tournee_id_produit_tournee_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.produit_tournee_id_produit_tournee_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.produit_tournee_id_produit_tournee_seq OWNER TO postgres;

--
-- Name: produit_tournee_id_produit_tournee_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.produit_tournee_id_produit_tournee_seq OWNED BY public.produit_tournee.id_produit_tournee;


--
-- Name: role; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.role (
    id_role integer NOT NULL,
    nom character varying(50) NOT NULL
);


ALTER TABLE public.role OWNER TO postgres;

--
-- Name: role_id_role_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.role_id_role_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.role_id_role_seq OWNER TO postgres;

--
-- Name: role_id_role_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.role_id_role_seq OWNED BY public.role.id_role;


--
-- Name: service; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service (
    id_service integer NOT NULL,
    nom character varying(100) NOT NULL,
    description text,
    id_competence integer NOT NULL,
    statut character varying(20) DEFAULT 'ACTIF'::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    id_agence integer NOT NULL,
    CONSTRAINT service_statut_check CHECK (((statut)::text = ANY ((ARRAY['ACTIF'::character varying, 'INACTIF'::character varying])::text[])))
);


ALTER TABLE public.service OWNER TO postgres;

--
-- Name: service_competence; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_competence (
    id_service integer NOT NULL,
    id_competence integer NOT NULL
);


ALTER TABLE public.service_competence OWNER TO postgres;

--
-- Name: service_id_service_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.service_id_service_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.service_id_service_seq OWNER TO postgres;

--
-- Name: service_id_service_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.service_id_service_seq OWNED BY public.service.id_service;


--
-- Name: stock; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock (
    id_stock integer NOT NULL,
    id_agence integer NOT NULL,
    quantite_disponible numeric(10,2) NOT NULL,
    date_entree timestamp without time zone NOT NULL
);


ALTER TABLE public.stock OWNER TO postgres;

--
-- Name: stock_id_stock_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.stock_id_stock_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.stock_id_stock_seq OWNER TO postgres;

--
-- Name: stock_id_stock_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.stock_id_stock_seq OWNED BY public.stock.id_stock;


--
-- Name: tournee; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tournee (
    id_tournee integer NOT NULL,
    id_destinataire integer NOT NULL,
    id_agence integer NOT NULL,
    id_benevole integer,
    id_planning integer,
    date_tournee date NOT NULL,
    statut public.statut_tournee DEFAULT 'en_attente'::public.statut_tournee NOT NULL,
    pdf_recapitulatif character varying(255)
);


ALTER TABLE public.tournee OWNER TO postgres;

--
-- Name: tournee_id_tournee_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tournee_id_tournee_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tournee_id_tournee_seq OWNER TO postgres;

--
-- Name: tournee_id_tournee_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tournee_id_tournee_seq OWNED BY public.tournee.id_tournee;


--
-- Name: utilisateur; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.utilisateur (
    id_utilisateur integer NOT NULL,
    nom character varying(100) NOT NULL,
    prenom character varying(100) NOT NULL,
    email character varying(150) NOT NULL,
    mot_de_passe character varying(255) NOT NULL,
    telephone character varying(20),
    adresse character varying(255),
    ville character varying(100),
    code_postal character varying(10),
    pays character varying(100),
    date_creation timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    dernier_login timestamp without time zone,
    email_verifie boolean DEFAULT false,
    etat_compte character varying(20) DEFAULT 'ACTIF'::character varying,
    id_role integer NOT NULL,
    id_agence integer,
    CONSTRAINT utilisateur_etat_compte_check CHECK (((etat_compte)::text = ANY ((ARRAY['ACTIF'::character varying, 'INACTIF'::character varying, 'EN_ATTENTE'::character varying])::text[])))
);


ALTER TABLE public.utilisateur OWNER TO postgres;

--
-- Name: utilisateur_id_utilisateur_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.utilisateur_id_utilisateur_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.utilisateur_id_utilisateur_seq OWNER TO postgres;

--
-- Name: utilisateur_id_utilisateur_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.utilisateur_id_utilisateur_seq OWNED BY public.utilisateur.id_utilisateur;


--
-- Name: adherent id_adherent; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.adherent ALTER COLUMN id_adherent SET DEFAULT nextval('public.adherent_id_adherent_seq'::regclass);


--
-- Name: agence id_agence; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.agence ALTER COLUMN id_agence SET DEFAULT nextval('public.agence_id_agence_seq'::regclass);


--
-- Name: association_beneficiaire id_association; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.association_beneficiaire ALTER COLUMN id_association SET DEFAULT nextval('public.association_beneficiaire_id_association_seq'::regclass);


--
-- Name: benevole id_benevole; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.benevole ALTER COLUMN id_benevole SET DEFAULT nextval('public.benevole_id_benevole_seq'::regclass);


--
-- Name: collecte id_collecte; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collecte ALTER COLUMN id_collecte SET DEFAULT nextval('public.collecte_id_collecte_seq'::regclass);


--
-- Name: commercant id_commercant; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.commercant ALTER COLUMN id_commercant SET DEFAULT nextval('public.commercant_id_commercant_seq'::regclass);


--
-- Name: competence id_competence; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.competence ALTER COLUMN id_competence SET DEFAULT nextval('public.competence_id_competence_seq'::regclass);


--
-- Name: cotisation id_cotisation; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cotisation ALTER COLUMN id_cotisation SET DEFAULT nextval('public.cotisation_id_cotisation_seq'::regclass);


--
-- Name: demande_service id_demande_service; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.demande_service ALTER COLUMN id_demande_service SET DEFAULT nextval('public.demande_service_id_demande_service_seq'::regclass);


--
-- Name: destinataire id_destinataire; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.destinataire ALTER COLUMN id_destinataire SET DEFAULT nextval('public.destinataire_id_destinataire_seq'::regclass);


--
-- Name: disponibilite id_disponibilite; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.disponibilite ALTER COLUMN id_disponibilite SET DEFAULT nextval('public.disponibilite_id_disponibilite_seq'::regclass);


--
-- Name: planning id_planning; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.planning ALTER COLUMN id_planning SET DEFAULT nextval('public.planning_id_planning_seq'::regclass);


--
-- Name: planning_excel id_planning_excel; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.planning_excel ALTER COLUMN id_planning_excel SET DEFAULT nextval('public.planning_excel_id_planning_excel_seq'::regclass);


--
-- Name: produit_collecte id_produit_collecte; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.produit_collecte ALTER COLUMN id_produit_collecte SET DEFAULT nextval('public.produit_collecte_id_produit_collecte_seq'::regclass);


--
-- Name: produit_tournee id_produit_tournee; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.produit_tournee ALTER COLUMN id_produit_tournee SET DEFAULT nextval('public.produit_tournee_id_produit_tournee_seq'::regclass);


--
-- Name: role id_role; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role ALTER COLUMN id_role SET DEFAULT nextval('public.role_id_role_seq'::regclass);


--
-- Name: service id_service; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service ALTER COLUMN id_service SET DEFAULT nextval('public.service_id_service_seq'::regclass);


--
-- Name: stock id_stock; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock ALTER COLUMN id_stock SET DEFAULT nextval('public.stock_id_stock_seq'::regclass);


--
-- Name: tournee id_tournee; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tournee ALTER COLUMN id_tournee SET DEFAULT nextval('public.tournee_id_tournee_seq'::regclass);


--
-- Name: utilisateur id_utilisateur; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.utilisateur ALTER COLUMN id_utilisateur SET DEFAULT nextval('public.utilisateur_id_utilisateur_seq'::regclass);


--
-- Data for Name: adherent; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.adherent (id_adherent, id_utilisateur, date_adhesion, date_expiration, statut) FROM stdin;
1	13	2026-07-31	2027-07-31	ACTIF
2	14	\N	\N	EN_ATTENTE
3	29	2026-08-10	2027-08-10	ACTIF
\.


--
-- Data for Name: agence; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.agence (id_agence, nom, adresse, ville, code_postal, pays, telephone, email) FROM stdin;
1	Siège Paris	24 Rue du Faubourg Saint-Antoine	Paris	75012	France	0142385566	siege.paris@nomorewaste.org
2	Agence Nantes	8 Rue des Entrepôts	Nantes	44000	France	0240123456	nantes@nomorewaste.org
3	Agence Marseille	25 Avenue du Prado	Marseille	13006	France	0491234567	marseille@nomorewaste.org
4	Agence Limoges	4 Rue des Écoles	Limoges	87000	France	0555987654	limoges@nomorewaste.org
5	Agence Naples	Via dei Mille 45	Naples	80121	Italie	0039817654321	naples@nomorewaste.org
6	Agence Porto	Rua de Santa Catarina 120	Porto	4000-102	Portugal	00351223456789	porto@nomorewaste.org
7	Agence Dublin	15 Grafton Street	Dublin	D02 F529	Irlande	00353876543210	dublin@nomorewaste.org
\.


--
-- Data for Name: association_beneficiaire; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.association_beneficiaire (id_association, id_utilisateur, nom_responsable, nom_association, nombre_beneficiaires, type_association, date_validation, statut) FROM stdin;
\.


--
-- Data for Name: benevole; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.benevole (id_benevole, id_utilisateur, permis, disponibilite, statut, id_competence, justificatif, commentaire_validation) FROM stdin;
2	3	t	SEMAINE	VALIDE	2	\N	\N
4	15	f	\N	VALIDE	5	uploads\\15.pdf	\N
10	27	f	\N	VALIDE	13	uploads\\27.pdf	\N
1	1	t	SEMAINE	VALIDE	2	\N	\N
5	22	t	\N	VALIDE	1	uploads\\22.pdf	\N
7	24	f	\N	VALIDE	2	uploads\\24.pdf	\N
8	25	f	\N	VALIDE	14	uploads\\25.pdf	\N
9	26	f	\N	VALIDE	5	uploads\\26.pdf	\N
6	23	t	SEMAINE	VALIDE	1	uploads\\23.pdf	\N
11	30	t	\N	VALIDE	1	uploads\\30.pdf	\N
\.


--
-- Data for Name: benevole_competence; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.benevole_competence (id_benevole, id_competence) FROM stdin;
\.


--
-- Data for Name: collecte; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.collecte (id_collecte, id_commercant, id_agence, id_benevole, id_planning, date_collecte, statut, pdf_recapitulatif) FROM stdin;
2	3	1	6	16	2026-08-13	effectuee	stockage\\livraisons\\livraison_2.pdf
3	4	1	6	17	2026-08-14	effectuee	stockage\\livraisons\\livraison_3.pdf
1	3	1	11	18	2026-08-12	effectuee	stockage\\livraisons\\livraison_1.pdf
\.


--
-- Data for Name: commercant; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.commercant (id_commercant, id_utilisateur, nom_entreprise, type_commerce, numero_siret, date_adhesion, date_expiration, cotisation, statut) FROM stdin;
3	7	Ofaweb	SUPERMARCHE	12345678900012	\N	\N	\N	ACTIF
4	8	France-Affaires	RESTAURANT	12345678900015	\N	\N	\N	ACTIF
\.


--
-- Data for Name: competence; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.competence (id_competence, nom, description) FROM stdin;
1	Chauffeur	Conduite des véhicules
2	Cuisine	Préparation de repas
3	Bricolage	Petits travaux divers
4	Électricité	Travaux électriques
5	Plomberie	Travaux de plomberie
6	Réparation	Réparation d'objets
7	Gardiennage	Surveillance de biens
8	Conseil anti-gaspillage	Sensibilisation au gaspillage
9	Tri des produits	Tri des denrées
10	Gestion des stocks	Organisation des stocks
11	Collecte des invendus	Collecte chez les commerçants
12	Distribution alimentaire	Distribution des produits
13	Manutention	Déplacement de marchandises
14	Animation d'atelier	Animation de groupes
15	Accueil	Accueil des adhérents
\.


--
-- Data for Name: cotisation; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.cotisation (id_cotisation, id_adherent, montant, date_paiement, date_debut, date_expiration, statut, stripe_session_id) FROM stdin;
1	1	5.00	2026-07-31	2026-07-31	2027-07-31	PAYEE	cs_test_a1OG5fkoYa0CpT8WocsCM4A5bzWkzjByQ1M8Wqfb9EHTJK7YY0qv4mweRV
2	3	5.00	2026-08-10	2026-08-10	2027-08-10	PAYEE	cs_test_a1Ima2sYSPFBNiprBLJ33QunNHBrcD9DUu7reUyPpHXnNBtVUhjSUH9uVD
\.


--
-- Data for Name: demande_service; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.demande_service (id_demande_service, id_service, id_adherent, date_demande, statut, date_creation, id_benevole, id_planning) FROM stdin;
1	1	1	2026-08-01	ATTRIBUE	2026-08-01 13:17:27.285462	2	10
2	1	3	2026-08-10	ATTRIBUE	2026-08-10 01:21:35.271038	7	14
\.


--
-- Data for Name: destinataire; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.destinataire (id_destinataire, id_agence, type, nom, adresse) FROM stdin;
1	1	ASSOCIATION	Nyota Raising Foundation	6 rue Charles Peguy
\.


--
-- Data for Name: disponibilite; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.disponibilite (id_disponibilite, id_benevole, date_disponibilite, heure_debut, heure_fin, statut, created_at) FROM stdin;
13	1	2026-08-14	11:00:00	13:00:00	DISPONIBLE	2026-08-02 18:55:48.288613
15	1	2026-08-17	14:00:00	16:00:00	DISPONIBLE	2026-08-02 18:56:35.918416
16	1	2026-08-18	15:00:00	16:00:00	DISPONIBLE	2026-08-02 18:56:56.266992
17	1	2026-08-19	10:00:00	13:00:00	DISPONIBLE	2026-08-02 18:57:48.411271
10	2	2026-08-11	13:00:00	15:00:00	RESERVE	2026-08-02 18:52:55.851343
9	2	2026-08-10	14:00:00	15:00:00	RESERVE	2026-08-02 18:52:07.013098
8	2	2026-08-07	10:00:00	13:00:00	RESERVE	2026-08-02 18:51:35.917207
4	2	2026-07-24	15:21:00	17:21:00	DISPONIBLE	2026-07-22 23:21:33.851461
3	2	2026-07-31	13:10:00	14:10:00	DISPONIBLE	2026-07-21 18:10:51.995251
6	2	2026-08-05	15:50:00	16:50:00	RESERVE	2026-08-02 18:50:09.377261
7	2	2026-08-06	12:00:00	14:00:00	RESERVE	2026-08-02 18:50:51.833713
11	1	2026-08-12	15:55:00	17:55:00	RESERVE	2026-08-02 18:55:04.985544
12	1	2026-08-13	14:00:00	15:00:00	RESERVE	2026-08-02 18:55:26.00845
18	7	2026-08-10	12:30:00	15:30:00	RESERVE	2026-08-10 00:27:10.95698
19	5	2026-08-12	11:00:00	13:00:00	RESERVE	2026-08-11 22:08:00.745804
21	6	2026-08-13	14:15:00	15:15:00	DISPONIBLE	2026-08-13 10:12:55.213514
20	6	2026-08-13	10:15:00	12:15:00	RESERVE	2026-08-13 10:12:28.846198
22	6	2026-08-14	10:20:00	12:20:00	RESERVE	2026-08-13 10:13:23.50874
23	11	2026-08-14	16:15:00	17:15:00	RESERVE	2026-08-14 15:12:10.118843
24	11	2026-08-14	17:20:00	18:20:00	RESERVE	2026-08-14 15:12:48.736964
25	11	2026-08-16	14:15:00	16:15:00	DISPONIBLE	2026-08-16 21:15:06.489182
26	2	2026-08-17	12:20:00	15:20:00	DISPONIBLE	2026-08-16 21:19:30.526863
27	11	2026-08-16	10:30:00	12:25:00	DISPONIBLE	2026-08-16 21:20:23.864534
28	11	2026-08-17	10:20:00	12:20:00	RESERVE	2026-08-16 21:21:30.105426
\.


--
-- Data for Name: planning; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.planning (id_planning, id_benevole, id_disponibilite, date, heure_debut, heure_fin, statut) FROM stdin;
10	2	6	2026-08-05	15:50:00	16:50:00	ATTRIBUE
11	2	7	2026-08-06	12:00:00	14:00:00	TERMINE
6	2	8	2026-08-07	10:00:00	13:00:00	TERMINE
14	7	18	2026-08-10	12:30:00	15:30:00	ATTRIBUE
5	2	9	2026-08-10	14:00:00	15:00:00	TERMINE
4	2	10	2026-08-11	13:00:00	15:00:00	TERMINE
15	5	19	2026-08-12	11:00:00	13:00:00	ATTRIBUE
12	1	11	2026-08-12	15:55:00	17:55:00	TERMINE
16	6	20	2026-08-13	10:15:00	12:15:00	ATTRIBUE
13	1	12	2026-08-13	14:00:00	15:00:00	TERMINE
17	6	22	2026-08-14	10:20:00	12:20:00	ATTRIBUE
18	11	23	2026-08-14	16:15:00	17:15:00	ATTRIBUE
19	11	24	2026-08-14	17:20:00	18:20:00	TERMINE
20	11	28	2026-08-17	10:20:00	12:20:00	ATTRIBUE
\.


--
-- Data for Name: planning_excel; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.planning_excel (id_planning_excel, id_planning, id_benevole, chemin_fichier, date_generation) FROM stdin;
1	14	7	stockage\\plannings\\planning_14.xlsx	2026-08-10 01:28:41.218273
2	16	6	stockage\\plannings\\planning_collecte_16.xlsx	2026-08-13 10:28:24.521507
3	17	6	stockage\\plannings\\planning_collecte_17.xlsx	2026-08-13 23:30:47.677077
4	18	11	stockage\\plannings\\planning_collecte_18.xlsx	2026-08-14 15:16:04.510474
\.


--
-- Data for Name: produit_collecte; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.produit_collecte (id_produit_collecte, id_collecte, libelle, quantite, code_barre, id_stock) FROM stdin;
4	2	Tomates	20.00	\N	\N
5	2	Pain	10.00	\N	\N
6	2	Yaourt	20.00	\N	\N
7	3	Bananes	30.00	NMW-7	1
8	3	Lait	20.00	NMW-8	2
1	1	Jus	20.00	NMW-1	3
2	1	Boite de conserve	30.00	NMW-2	4
3	1	Pomme	50.00	NMW-3	5
\.


--
-- Data for Name: produit_tournee; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.produit_tournee (id_produit_tournee, id_tournee, id_stock, quantite) FROM stdin;
1	1	5	50.00
\.


--
-- Data for Name: role; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.role (id_role, nom) FROM stdin;
1	ADMIN_GENERAL
2	ADMIN_AGENCE
3	COMMERCANT
4	BENEVOLE
5	ADHERENT
6	ASSOCIATION
\.


--
-- Data for Name: service; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.service (id_service, nom, description, id_competence, statut, created_at, id_agence) FROM stdin;
1	Cours de cuisine	Cours où on apprend à bien cuisiner	2	ACTIF	2026-07-31 22:55:04.43556	1
3	Plomberie	Tavaux de plomberie	5	ACTIF	2026-08-09 22:22:14.314103	1
4	Bricolage	Bricolez ce que vous voulez	3	ACTIF	2026-08-09 22:22:41.21294	1
6	Conseil anti-gaspillage	Conseils d'expert	8	ACTIF	2026-08-09 22:24:34.180748	1
7	Gardiennage	Gardiannage à domicile	7	ACTIF	2026-08-09 22:25:05.914584	1
8	Manutention	Chargement et déchargement de colis	13	ACTIF	2026-08-09 22:26:39.08649	1
9	Bricolage	Bricolage de tout type de meuble	3	ACTIF	2026-08-09 22:29:54.20306	1
10	Conduite	Faites vous-conduire	1	ACTIF	2026-08-09 22:30:40.608691	1
11	Animation d'atelier	Animer tous vos ateliers	14	ACTIF	2026-08-09 22:31:23.215157	1
5	Réparation meuble	Ce que vous voulez	6	ACTIF	2026-08-09 22:23:43.097439	1
\.


--
-- Data for Name: service_competence; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.service_competence (id_service, id_competence) FROM stdin;
\.


--
-- Data for Name: stock; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stock (id_stock, id_agence, quantite_disponible, date_entree) FROM stdin;
1	1	30.00	2026-08-13 23:31:24.591495
2	1	20.00	2026-08-13 23:31:24.591495
3	1	20.00	2026-08-14 15:16:52.787518
4	1	30.00	2026-08-14 15:16:52.787518
5	1	0.00	2026-08-14 15:16:52.787518
\.


--
-- Data for Name: tournee; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tournee (id_tournee, id_destinataire, id_agence, id_benevole, id_planning, date_tournee, statut, pdf_recapitulatif) FROM stdin;
1	1	1	11	20	2026-08-15	effectuee	stockage\\livraisons\\livraison_tournee_1.pdf
\.


--
-- Data for Name: utilisateur; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.utilisateur (id_utilisateur, nom, prenom, email, mot_de_passe, telephone, adresse, ville, code_postal, pays, date_creation, dernier_login, email_verifie, etat_compte, id_role, id_agence) FROM stdin;
27	Bentala	Karim	kbentala@gmail.com	$2a$10$669fzQEu.nn1aRIAWRyUO.rDJLt05n29zxQkJhoZfK31JtUrWL2Fq	06100006	12 rue des Lieux	Paris	75007	France	2026-08-09 23:49:54.418712	\N	f	ACTIF	4	1
14	Lupaka	Nathan	nlupaka@gmail.com	$2a$10$9eINyrkaUBw1PUHagOPriuW4APv6aGBkOa71LP6QyTlIANnbFMXo2	+33758537264	6 rue Mallet	Paris	75017	France	2026-07-31 11:37:55.229144	\N	f	ACTIF	5	1
2	admin	admin	adminprincipal@gmail.com	$2a$10$yS/.lzQP8G/JskguoUukrOej6EsghFSDtsdZN3O/8km7A2v8YHRXi						2026-07-12 20:53:31.409902	\N	f	ACTIF	1	\N
15	motsepe	patrick	pmotsepe@gmail.com	$2a$10$KdtIwjsG2QdbbD4bumMaOubO8kZipZJwCeATPLVUgWkNWzSdP363u	+33634043801	7 avenue de l'étoile	Paris	75001	France	2026-08-03 11:17:43.218058	\N	f	ACTIF	4	1
13	Zola	Nato	nzola@gmail.com	$2a$10$KmCGu6DIAvUizfg0cZ97Ve8TLCZoOmF04RdgCV9liA7Wjc/oM5C26	+33634033701	4 rue Charles Peguy	Paris	75010	France	2026-07-30 10:57:07.557535	\N	f	ACTIF	5	1
7	Laroura	David	davidlaroura@gmail.com	$2a$10$VIusDl0A4VXLOj6cOTZZrOOwFAJhXxsxqlJEg2BXktcyTZFLHJvlS	+33754256398	24 rue du Ballon 	Paris	93100	France	2026-07-16 21:55:55.71421	\N	f	ACTIF	3	1
60	test	test	test@gmail.com	root	+3378945612	azertyuipo	Paris	de	ezf	2026-07-14 18:30:56.151894	\N	f	ACTIF	4	\N
10	adminParis	adminParis	adminparis@gmail.com	$2a$10$HjDqP.i8yilu4IXDRt8enu45ZUMjJm52nECP/5TlnfyCjw3W.C75S	+33674123654	5 rue Maillet	Paris	75001	France	2026-07-19 00:10:36.723941	\N	f	ACTIF	2	1
1	Becker	Alain	alainbecker@gmail.com	$2a$10$FXMo1./jGOEj8ypoZEjsE.lO6FsN2ZGw2pD6PmHJjSKzdJJBE/o1i	+33634033701	Paris	Paris		France	2026-07-08 23:29:11.184896	\N	f	ACTIF	4	1
29	Bawhere	Trésor	tbawhere@gmail.com	$2a$10$3RjuK2qslEfnUBJiGRGwNOw6Dt7GOnoFSDS25Wp/3sCzvNCIXkIqO	06100007	12 rue des Lieux	Paris	75013	France	2026-08-10 00:42:33.140837	\N	f	ACTIF	5	1
16	adminNantes	adminNantes	adminnantes@gmail.com	$2a$10$Z0sF9D8tM94mtEX/8ZH27edS0zBrSllyAW3gMStuEdK313vBcP/s6	+33754256398	3 rue Aimé Césaire	Nantes	93400	France	2026-08-09 23:11:19.37565	\N	f	ACTIF	2	2
17	adminMarseilles	adminMarseilles	adminmarseilles@gmail.com	$2a$10$r7pZUS.hx019vG1BqTDgMeHqOv0.WPAkRW/8hq8iRGUQ5khhHyI4W	0758537267	3 rue Aimé Césaire	Marseilles	93400	France	2026-08-09 23:13:09.264885	\N	f	ACTIF	2	3
18	adminLimoges	adminLimoges	adminlimoges@gmail.com	$2a$10$tc9Z1dyjI7FMNW.oNKKrVe/MFo5eD1F1woFcWQT0uPOH9rPJnNKFu	0634033701	3 rue Aimé Césaire	Limoges	93400	France	2026-08-09 23:14:29.668968	\N	f	ACTIF	2	4
19	adminNaples	adminNaples	adminnaples@gmail.com	$2a$10$FfhGCFEuQ4MnLR8RdG/8Aesr9s5eQtYSEAdBf4BZ/xgA5c8bgoyyG	0758537265	3 rue Aimé Césaire	Naples	93400	Italie	2026-08-09 23:15:18.011857	\N	f	ACTIF	2	5
20	adminPorto	adminPorto	adminporto@gmail.com	$2a$10$BiT9cnnG5uXSk0JM8AWvfOr63EQnQUemFvJJk91xuTzKsgiCgoXqi	0634033707	3 rue Aimé Césaire	Porto	93400	Portugal	2026-08-09 23:20:52.088809	\N	f	ACTIF	2	6
21	adminDublin	adminDublin	admindublin@gmail.com	$2a$10$Tw.wiMYRJgMvFAOglGnp4uCkyEAJ74YcSXYvVIC.WKieNkB9gDK2S	0634033774	3 rue Aimé Césaire	Dublin	93400	Irlande	2026-08-09 23:21:38.223025	\N	f	ACTIF	2	7
3	Dupont	Gerard	gdupont@gmail.com	$2a$10$9lpD9we4T06ef5hODzoiFeWIdPeHbLP64VSinuWqHfEYjSG3NdsHi	+3345236987 	12 rue de Conflants	Paris	75015	France	2026-07-12 22:40:26.658396	\N	f	ACTIF	4	1
8	Adelaid	Olivier	olivieradelaid@gmail.com	$2a$10$GrX.P0CSe/LwM61x/5jMT.XBWapG.WIrWmF9Gmy49SH2f7BOdIMXy	+33674123658	6 rue Adolphe Focillon	Paris	75001	France	2026-07-16 23:09:34.034205	\N	f	ACTIF	3	1
24	Bernard	Sophie 	sbernard@gmail.com	$2a$10$/bq5W6KRI/A4zlN61/MHuO3w1XQ5fN.XAEqge1MPNFzuRAcjYpDji	061000003	12 rue de Conflants	Paris	75003	France	2026-08-09 23:41:36.482255	\N	f	ACTIF	4	1
25	Girard	Camille	cgirard@gmail.com	$2a$10$iwwz6qei1EbvxRfkTi3Iy.ySjbsrEdCVkU.1lIoVQ4ey2CTllYX76	06100004	12 rue de Conflants	Paris	75005	France	2026-08-09 23:45:29.97809	\N	f	ACTIF	4	1
26	Roux	Elodie 	eroux@gmail.com	$2a$10$V4A.zi9fLAhPp2hHNFzddufXCHxlpJmD0KpY7a/dh9tBlMYydPKZy	06100005	12 rue de Conflants	Paris	75006	France	2026-08-09 23:47:09.125678	\N	f	ACTIF	4	1
22	Lefevre	Marc	mlefevre@gmail.com	$2a$10$HY/uxlmaq29Q.d.SfJ2KQuZk1PGapAHzlWb4VH6px/56h50.0ByFu	0611000001	12 rue des Lieux	Paris	75011	France	2026-08-09 23:37:56.592755	\N	f	ACTIF	4	1
23	Rousseau	Antoine	arousseau@gmail.com	$2a$10$2mPF/BJP4Mv7.pLLbOOQO.d4pZU0U/508jTI.P0rpczUDgPfUH7xG	06100002	12 rue des Lieux	Paris		France	2026-08-09 23:40:17.691038	\N	f	ACTIF	4	1
30	Okita	Teddy	tokita@gmail.com	$2a$10$ZxHbdDbLfBuKGi6ja.tBceQZmhIO1h5pYh37bT4VHAGN1qe1Nrhma	+3345236984	4 rue Charles Peguy	Créteil	9400	France	2026-08-14 15:08:07.807842	\N	f	ACTIF	4	1
\.


--
-- Name: adherent_id_adherent_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.adherent_id_adherent_seq', 3, true);


--
-- Name: agence_id_agence_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.agence_id_agence_seq', 14, true);


--
-- Name: association_beneficiaire_id_association_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.association_beneficiaire_id_association_seq', 1, false);


--
-- Name: benevole_id_benevole_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.benevole_id_benevole_seq', 11, true);


--
-- Name: collecte_id_collecte_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.collecte_id_collecte_seq', 3, true);


--
-- Name: commercant_id_commercant_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.commercant_id_commercant_seq', 5, true);


--
-- Name: competence_id_competence_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.competence_id_competence_seq', 15, true);


--
-- Name: cotisation_id_cotisation_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.cotisation_id_cotisation_seq', 2, true);


--
-- Name: demande_service_id_demande_service_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.demande_service_id_demande_service_seq', 2, true);


--
-- Name: destinataire_id_destinataire_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.destinataire_id_destinataire_seq', 1, true);


--
-- Name: disponibilite_id_disponibilite_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.disponibilite_id_disponibilite_seq', 28, true);


--
-- Name: planning_excel_id_planning_excel_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.planning_excel_id_planning_excel_seq', 4, true);


--
-- Name: planning_id_planning_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.planning_id_planning_seq', 20, true);


--
-- Name: produit_collecte_id_produit_collecte_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.produit_collecte_id_produit_collecte_seq', 8, true);


--
-- Name: produit_tournee_id_produit_tournee_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.produit_tournee_id_produit_tournee_seq', 1, true);


--
-- Name: role_id_role_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.role_id_role_seq', 6, true);


--
-- Name: service_id_service_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.service_id_service_seq', 11, true);


--
-- Name: stock_id_stock_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.stock_id_stock_seq', 5, true);


--
-- Name: tournee_id_tournee_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.tournee_id_tournee_seq', 1, true);


--
-- Name: utilisateur_id_utilisateur_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.utilisateur_id_utilisateur_seq', 30, true);


--
-- Name: adherent adherent_id_utilisateur_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.adherent
    ADD CONSTRAINT adherent_id_utilisateur_key UNIQUE (id_utilisateur);


--
-- Name: adherent adherent_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.adherent
    ADD CONSTRAINT adherent_pkey PRIMARY KEY (id_adherent);


--
-- Name: agence agence_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.agence
    ADD CONSTRAINT agence_pkey PRIMARY KEY (id_agence);


--
-- Name: association_beneficiaire association_beneficiaire_id_utilisateur_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.association_beneficiaire
    ADD CONSTRAINT association_beneficiaire_id_utilisateur_key UNIQUE (id_utilisateur);


--
-- Name: association_beneficiaire association_beneficiaire_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.association_beneficiaire
    ADD CONSTRAINT association_beneficiaire_pkey PRIMARY KEY (id_association);


--
-- Name: benevole_competence benevole_competence_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.benevole_competence
    ADD CONSTRAINT benevole_competence_pkey PRIMARY KEY (id_benevole, id_competence);


--
-- Name: benevole benevole_id_utilisateur_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.benevole
    ADD CONSTRAINT benevole_id_utilisateur_key UNIQUE (id_utilisateur);


--
-- Name: benevole benevole_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.benevole
    ADD CONSTRAINT benevole_pkey PRIMARY KEY (id_benevole);


--
-- Name: collecte collecte_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collecte
    ADD CONSTRAINT collecte_pkey PRIMARY KEY (id_collecte);


--
-- Name: commercant commercant_id_utilisateur_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.commercant
    ADD CONSTRAINT commercant_id_utilisateur_key UNIQUE (id_utilisateur);


--
-- Name: commercant commercant_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.commercant
    ADD CONSTRAINT commercant_pkey PRIMARY KEY (id_commercant);


--
-- Name: competence competence_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.competence
    ADD CONSTRAINT competence_pkey PRIMARY KEY (id_competence);


--
-- Name: cotisation cotisation_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cotisation
    ADD CONSTRAINT cotisation_pkey PRIMARY KEY (id_cotisation);


--
-- Name: cotisation cotisation_stripe_session_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cotisation
    ADD CONSTRAINT cotisation_stripe_session_id_key UNIQUE (stripe_session_id);


--
-- Name: demande_service demande_service_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.demande_service
    ADD CONSTRAINT demande_service_pkey PRIMARY KEY (id_demande_service);


--
-- Name: destinataire destinataire_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.destinataire
    ADD CONSTRAINT destinataire_pkey PRIMARY KEY (id_destinataire);


--
-- Name: disponibilite disponibilite_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.disponibilite
    ADD CONSTRAINT disponibilite_pkey PRIMARY KEY (id_disponibilite);


--
-- Name: planning_excel planning_excel_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.planning_excel
    ADD CONSTRAINT planning_excel_pkey PRIMARY KEY (id_planning_excel);


--
-- Name: planning planning_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.planning
    ADD CONSTRAINT planning_pkey PRIMARY KEY (id_planning);


--
-- Name: produit_collecte produit_collecte_code_barre_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.produit_collecte
    ADD CONSTRAINT produit_collecte_code_barre_key UNIQUE (code_barre);


--
-- Name: produit_collecte produit_collecte_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.produit_collecte
    ADD CONSTRAINT produit_collecte_pkey PRIMARY KEY (id_produit_collecte);


--
-- Name: produit_tournee produit_tournee_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.produit_tournee
    ADD CONSTRAINT produit_tournee_pkey PRIMARY KEY (id_produit_tournee);


--
-- Name: role role_nom_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role
    ADD CONSTRAINT role_nom_key UNIQUE (nom);


--
-- Name: role role_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role
    ADD CONSTRAINT role_pkey PRIMARY KEY (id_role);


--
-- Name: service_competence service_competence_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_competence
    ADD CONSTRAINT service_competence_pkey PRIMARY KEY (id_service, id_competence);


--
-- Name: service service_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service
    ADD CONSTRAINT service_pkey PRIMARY KEY (id_service);


--
-- Name: stock stock_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock
    ADD CONSTRAINT stock_pkey PRIMARY KEY (id_stock);


--
-- Name: tournee tournee_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tournee
    ADD CONSTRAINT tournee_pkey PRIMARY KEY (id_tournee);


--
-- Name: utilisateur utilisateur_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.utilisateur
    ADD CONSTRAINT utilisateur_email_key UNIQUE (email);


--
-- Name: utilisateur utilisateur_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.utilisateur
    ADD CONSTRAINT utilisateur_pkey PRIMARY KEY (id_utilisateur);


--
-- Name: benevole benevole_id_competence_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.benevole
    ADD CONSTRAINT benevole_id_competence_fkey FOREIGN KEY (id_competence) REFERENCES public.competence(id_competence);


--
-- Name: collecte collecte_id_agence_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collecte
    ADD CONSTRAINT collecte_id_agence_fkey FOREIGN KEY (id_agence) REFERENCES public.agence(id_agence);


--
-- Name: collecte collecte_id_benevole_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collecte
    ADD CONSTRAINT collecte_id_benevole_fkey FOREIGN KEY (id_benevole) REFERENCES public.benevole(id_benevole);


--
-- Name: collecte collecte_id_commercant_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collecte
    ADD CONSTRAINT collecte_id_commercant_fkey FOREIGN KEY (id_commercant) REFERENCES public.commercant(id_commercant);


--
-- Name: collecte collecte_id_planning_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collecte
    ADD CONSTRAINT collecte_id_planning_fkey FOREIGN KEY (id_planning) REFERENCES public.planning(id_planning);


--
-- Name: demande_service demande_service_id_benevole_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.demande_service
    ADD CONSTRAINT demande_service_id_benevole_fkey FOREIGN KEY (id_benevole) REFERENCES public.benevole(id_benevole);


--
-- Name: demande_service demande_service_id_planning_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.demande_service
    ADD CONSTRAINT demande_service_id_planning_fkey FOREIGN KEY (id_planning) REFERENCES public.planning(id_planning);


--
-- Name: destinataire destinataire_id_agence_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.destinataire
    ADD CONSTRAINT destinataire_id_agence_fkey FOREIGN KEY (id_agence) REFERENCES public.agence(id_agence);


--
-- Name: cotisation fk_adherent_cotisation; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cotisation
    ADD CONSTRAINT fk_adherent_cotisation FOREIGN KEY (id_adherent) REFERENCES public.adherent(id_adherent);


--
-- Name: utilisateur fk_agence; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.utilisateur
    ADD CONSTRAINT fk_agence FOREIGN KEY (id_agence) REFERENCES public.agence(id_agence);


--
-- Name: benevole_competence fk_benevole; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.benevole_competence
    ADD CONSTRAINT fk_benevole FOREIGN KEY (id_benevole) REFERENCES public.benevole(id_benevole);


--
-- Name: planning fk_benevole; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.planning
    ADD CONSTRAINT fk_benevole FOREIGN KEY (id_benevole) REFERENCES public.benevole(id_benevole);


--
-- Name: benevole_competence fk_competence; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.benevole_competence
    ADD CONSTRAINT fk_competence FOREIGN KEY (id_competence) REFERENCES public.competence(id_competence);


--
-- Name: demande_service fk_demande_adherent; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.demande_service
    ADD CONSTRAINT fk_demande_adherent FOREIGN KEY (id_adherent) REFERENCES public.adherent(id_adherent) ON DELETE CASCADE;


--
-- Name: demande_service fk_demande_service; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.demande_service
    ADD CONSTRAINT fk_demande_service FOREIGN KEY (id_service) REFERENCES public.service(id_service) ON DELETE RESTRICT;


--
-- Name: planning fk_disponibilite; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.planning
    ADD CONSTRAINT fk_disponibilite FOREIGN KEY (id_disponibilite) REFERENCES public.disponibilite(id_disponibilite);


--
-- Name: disponibilite fk_disponibilite_benevole; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.disponibilite
    ADD CONSTRAINT fk_disponibilite_benevole FOREIGN KEY (id_benevole) REFERENCES public.benevole(id_benevole) ON DELETE CASCADE;


--
-- Name: utilisateur fk_role; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.utilisateur
    ADD CONSTRAINT fk_role FOREIGN KEY (id_role) REFERENCES public.role(id_role);


--
-- Name: adherent fk_utilisateur_adherent; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.adherent
    ADD CONSTRAINT fk_utilisateur_adherent FOREIGN KEY (id_utilisateur) REFERENCES public.utilisateur(id_utilisateur);


--
-- Name: association_beneficiaire fk_utilisateur_association; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.association_beneficiaire
    ADD CONSTRAINT fk_utilisateur_association FOREIGN KEY (id_utilisateur) REFERENCES public.utilisateur(id_utilisateur);


--
-- Name: benevole fk_utilisateur_benevole; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.benevole
    ADD CONSTRAINT fk_utilisateur_benevole FOREIGN KEY (id_utilisateur) REFERENCES public.utilisateur(id_utilisateur) ON DELETE CASCADE;


--
-- Name: commercant fk_utilisateur_commercant; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.commercant
    ADD CONSTRAINT fk_utilisateur_commercant FOREIGN KEY (id_utilisateur) REFERENCES public.utilisateur(id_utilisateur);


--
-- Name: planning_excel planning_excel_id_benevole_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.planning_excel
    ADD CONSTRAINT planning_excel_id_benevole_fkey FOREIGN KEY (id_benevole) REFERENCES public.benevole(id_benevole);


--
-- Name: planning_excel planning_excel_id_planning_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.planning_excel
    ADD CONSTRAINT planning_excel_id_planning_fkey FOREIGN KEY (id_planning) REFERENCES public.planning(id_planning);


--
-- Name: produit_collecte produit_collecte_id_collecte_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.produit_collecte
    ADD CONSTRAINT produit_collecte_id_collecte_fkey FOREIGN KEY (id_collecte) REFERENCES public.collecte(id_collecte);


--
-- Name: produit_collecte produit_collecte_id_stock_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.produit_collecte
    ADD CONSTRAINT produit_collecte_id_stock_fkey FOREIGN KEY (id_stock) REFERENCES public.stock(id_stock);


--
-- Name: produit_tournee produit_tournee_id_stock_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.produit_tournee
    ADD CONSTRAINT produit_tournee_id_stock_fkey FOREIGN KEY (id_stock) REFERENCES public.stock(id_stock);


--
-- Name: produit_tournee produit_tournee_id_tournee_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.produit_tournee
    ADD CONSTRAINT produit_tournee_id_tournee_fkey FOREIGN KEY (id_tournee) REFERENCES public.tournee(id_tournee);


--
-- Name: service_competence service_competence_id_competence_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_competence
    ADD CONSTRAINT service_competence_id_competence_fkey FOREIGN KEY (id_competence) REFERENCES public.competence(id_competence);


--
-- Name: service_competence service_competence_id_service_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_competence
    ADD CONSTRAINT service_competence_id_service_fkey FOREIGN KEY (id_service) REFERENCES public.service(id_service);


--
-- Name: service service_id_agence_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service
    ADD CONSTRAINT service_id_agence_fkey FOREIGN KEY (id_agence) REFERENCES public.agence(id_agence);


--
-- Name: service service_id_competence_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service
    ADD CONSTRAINT service_id_competence_fkey FOREIGN KEY (id_competence) REFERENCES public.competence(id_competence);


--
-- Name: stock stock_id_agence_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock
    ADD CONSTRAINT stock_id_agence_fkey FOREIGN KEY (id_agence) REFERENCES public.agence(id_agence);


--
-- Name: tournee tournee_id_agence_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tournee
    ADD CONSTRAINT tournee_id_agence_fkey FOREIGN KEY (id_agence) REFERENCES public.agence(id_agence);


--
-- Name: tournee tournee_id_benevole_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tournee
    ADD CONSTRAINT tournee_id_benevole_fkey FOREIGN KEY (id_benevole) REFERENCES public.benevole(id_benevole);


--
-- Name: tournee tournee_id_destinataire_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tournee
    ADD CONSTRAINT tournee_id_destinataire_fkey FOREIGN KEY (id_destinataire) REFERENCES public.destinataire(id_destinataire);


--
-- Name: tournee tournee_id_planning_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tournee
    ADD CONSTRAINT tournee_id_planning_fkey FOREIGN KEY (id_planning) REFERENCES public.planning(id_planning);


--
-- PostgreSQL database dump complete
--

\unrestrict RTsRvWeYoU4RHyP6dxElvXPuDz6ji7fjTT5Rb0koqeL8bIYvPWuZimWhqp5a0Bb

