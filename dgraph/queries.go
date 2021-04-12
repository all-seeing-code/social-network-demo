package dgraph

type input struct {
	query string
	vars  map[string]string
}

var queries = map[string]input{
	"IS01": {
		query: `query all($ID: string){
			q(func: eq(fqid, $ID)) {
				firstName
				lastName
				birthday
				locationIP
				browserUsed
				gender
				creationDate
				isLocatedIn{
					id
					name
				}
			}
		}`,
		vars: map[string]string{"$ID": "person_4398046514948"},
	},
	"IS03": {
		query: `query all($ID: string){
			q(func: eq(fqid, $ID)) {
				knows @facets( orderdesc: creationDate, orderdesc: id){
					id
					firstName
					lastName
				}
			}
		}`,
		vars: map[string]string{"$ID": "person_933"},
	},
	"IS06": {
		query: `query all($ID: string){
			q(func: eq(fqid, $ID)) {
				~containerOf{
					id
					title
					hasModerator{
						id
						firstName
						lastName
					}
				}
			}
		}`,
		vars: map[string]string{"$ID": "post_3"},
	},
	"IC06": {
		query: `query all($ID: int, $tagName: string){
			pid as var(func: type(person))@filter(eq(id, $ID))
			var(func: uid(pid)){
				f11 as knows{
					f12 as knows @filter(NOT uid(pid))
					f13 as ~knows @filter(NOT uid(pid))
				}
				f21 as ~knows{
					f22 as knows @filter(NOT uid(pid))
					f23 as ~knows @filter(NOT uid(pid))
				}
			}
			# tag is the tag that's passed in.
			tag as var(func: eq(name, $tagName)) {
				# posts are all the posts which have this tag.
				posts as ~hasTag @filter(type(post))
			}
			# Collect posts by friends and the tags of those posts.
			# Restrict posts to which have the tag
			var(func: uid(f11 ,f12, f13, f21, f22, f23)){
				tc as math(1) # variable propagation trick to count all paths from post -> tag.
				~hasCreator @filter(uid(posts)) {
					otherTag as hasTag {
					    # tagCount tracks tag -> post count, by using variable propagation.
						tagCount as math(tc)
					}
				}
			}
			tags as var(func: uid(otherTag)) @filter(NOT uid(tag))
			q(func: uid(tags), orderdesc: val(tagCount), first:10){
				name
				postCount: val(tagCount)
			}
		}`,
		vars: map[string]string{"$ID": "102", "$tagName": "Rumi"},
	},
	"IC07": {
		query: `query all($ID: int){
			pid as var(func: type(person))@filter(eq(id, $ID))
			var(func: uid(pid))@cascade {
				messages as ~hasCreator {
					personLikes as ~likes
				}
				friends1 as knows{
					friend1 as math(1)
				}
				friends2 as ~knows{
					friend2 as math(1)
				}
			}
			q(func: uid(friends1, friends2, personLikes), orderdesc: id, first:20) @filter(uid_in(likes, uid(messages) )){
				isFriend2: val(friend2)
				isFriend1: val(friend1)
				id
				firstName
				lastName
				likes @facets(orderdesc: creationDate) @filter(uid(messages)) (first:1){
					id
					creationDate
				}
			}
		}`,
		vars: map[string]string{"$ID": "102"},
	},
	"IC10": {
		query: `query all($ID: int){
			pid as var(func: type(person))@filter(eq(id, $ID))
			var(func: uid(pid)){
				exf1 as knows{
					f1 as knows
					f2 as ~knows
				}
				exf2 as ~knows{
					f3 as knows
					f4 as ~knows
				}
			}
			friendsOfInterest as var(func:uid(f1,f2,f3,f4)) @filter(NOT uid(exf1, exf2, pid))
			var(func: uid(pid)){
				tagsOfInterest as hasInterest (orderasc: fqid)
			}
			var(func: uid(friendsOfInterest)){
				common as count(~hasCreator) @filter(uid_in(hasTag, uid(tagsOfInterest)) AND type(post))
			}
			var(func: uid(friendsOfInterest)){
				uncommon as count(~hasCreator) @filter(NOT uid_in(hasTag, uid(tagsOfInterest)) AND type(post))
			}
			var(func: uid(friendsOfInterest) ){
				interest as math(common - uncommon )
			}
			q(func: uid(friendsOfInterest), orderdesc: val(interest), first:10){
				fqid
				firstName
				lastName
				gender
				isLocatedIn{
					name
				}
				co: val(common)
				un: val(uncommon)
				interest: val(interest)
			}
		}`,
		vars: map[string]string{"$ID": "933"},
	},
	"IC11": {
		query: `query all($ID: int, $name: string, $workFrom: string){
			pid as var(func: type(person))@filter(eq(id, $ID))
			var(func: uid(pid)){
				f1 as knows @filter(NOT uid(pid)){
					f11 as knows @filter(NOT uid(pid))
					f12 as ~knows @filter(NOT uid(pid))
				}
				f2 as ~knows @filter(NOT uid(pid)){
					f21 as knows @filter(NOT uid(pid))
					f22 as ~knows @filter(NOT uid(pid))
				}
			}
			friendsAndFoF as var(func:uid(f1, f2, f11, f12, f21, f22))
			country as var(func: eq(name, $name))
			organisation as var(func: type(organisation)) @filter(uid_in(isLocatedIn, uid(country)))
			relevantFriends as var(func: uid(friendsAndFoF))@filter(uid_in(workAt, uid(organisation)))
			q(func: uid(relevantFriends), orderasc: id) @cascade{
					id
					firstName
					lastName
					workAt @filter( uid(organisation)) @facets(workFrom) @facets(le(workFrom, $workFrom)){
						id
						name
					}
			}
		}`,
		// Put date from 2007 because < includes dates from 2008 as well.
		vars: map[string]string{"$ID": "4398046514948", "$name": "Sri_Lanka", "$workFrom": "2007-12-31T18:52:55.426+00:00"},
	},
	"IC13": {
		query: `query all($ID1: int, $ID2: int){
			A as var(func: eq(fqid, $ID1))
			M as var(func: eq(fqid, $ID2))
			path as shortest(from: uid(A), to: uid(M)) {
			 knows
			 ~knows
			}
			path(func: uid(path)) {
			  id
			}
		}`,
		vars: map[string]string{"$ID1": "398046514948", "$ID2": "28587302326035"},
	},
}
